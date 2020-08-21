package model

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"log"
	"math"
	"net/http"
	"prom2lyrid/logger"
	"sync"
	"time"
)

// exporter type list:
// windows
// node
// mongo
// unknown

// Endpoint Status Enums:
// Starting
// Started
// Warning
// Error
// Stopping
// Stopped
//

type ExporterEndpoint struct {
	ID           string       `json:"id"`
	Gateway      string       `json:gateway`
	URL          string       `json:"url"`
	Config       ScrapeConfig `json:"config"`
	ExporterType string       `json:"exportertype"`

	Status           string            `json:"status"`
	LastScrape       time.Time         `json:"last_scrape"`
	AdditionalLabels map[string]string `json:"additional_labels"`

	Message   string `json:"message"`
	IsUpdated bool   `json:"is_updated"`
	IsCompress bool  `json:"is_compress"`

	//IsScraping bool `json:-`
	//Stopping

	DurationSinceLastUpdate time.Duration
	LastUpdateTime          time.Time
	LastScrapeDuration      time.Duration

	mux    sync.Mutex
	Result []*dto.MetricFamily `json:"-"`
	ctx    context.Context
	cancel context.CancelFunc
}

func CreateEndpoint(url string) ExporterEndpoint {
	return ExporterEndpoint{
		ID:               uuid.New().String(),
		URL:              url,
		Config:           CreateDefaultScrapeConfig(),
		ExporterType:     "unknown",
		AdditionalLabels: map[string]string{},
		Status:           "Starting",
		IsCompress:		  true,
	}
}

func (endpoint *ExporterEndpoint) SetUpdate(update bool) {
	endpoint.mux.Lock()
	if update {
		endpoint.ResetTime()
	}
	endpoint.IsUpdated = update
	endpoint.mux.Unlock()
}

func (endpoint *ExporterEndpoint) SetTimeDuration() {
	endpoint.DurationSinceLastUpdate = time.Since(endpoint.LastUpdateTime)
}

func (endpoint *ExporterEndpoint) ResetTime() {
	endpoint.DurationSinceLastUpdate = 0
	endpoint.LastUpdateTime = time.Now().UTC()
}

func (endpoint *ExporterEndpoint) Scrape() ([]*dto.MetricFamily, error) {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	transport.ResponseHeaderTimeout = time.Minute

	mfChan := make(chan *dto.MetricFamily, 1024)

	err := prom2json.FetchMetricFamilies(endpoint.URL, mfChan, transport)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := []*dto.MetricFamily{}
	for mf := range mfChan {
		result = append(result, mf)
	}

	endpoint.LastScrape = time.Now().UTC()
	endpoint.ResetTime()

	return result, nil
}

func (endpoint *ExporterEndpoint) Run(ctx context.Context) {

	endpoint.ctx, endpoint.cancel = context.WithCancel(ctx)
	duration, _ := time.ParseDuration(endpoint.Config.ScrapeInterval)
	for c := time.Tick(duration); ; {

		if endpoint.Status == "Error" {
			// do not scrape
		} else {
			level.Info(logger.GetInstance().Logger).Log("Message", "Running endpoint", "Endpoint",  endpoint.URL)
			start := time.Now()
			result, err := endpoint.Scrape()
			scrapDuration := time.Now().Sub(start)
			if err == nil {
				endpoint.SetUpdate(true)
				endpoint.Status = "Running"
				endpoint.Message = ""
				endpoint.LastScrapeDuration = scrapDuration

				for _, metricfamily := range result {
					for _, metric := range metricfamily.Metric {
						if metric.Summary != nil {
							for _, quantile := range metric.Summary.Quantile {
								if quantile.Value != nil {
									if math.IsNaN(*quantile.Value) || math.IsInf(*quantile.Value, 1) || math.IsInf(*quantile.Value, -1) {
										quantile.Value = nil
									}
								}
							}
						}
					}
				}

				endpoint.Result = result
				level.Info(logger.GetInstance().Logger).Log("Endpoint",  endpoint.URL, "ScrapeTime(ms)", scrapDuration.Milliseconds())
			} else {
				level.Error(logger.GetInstance().Logger).Log("Message",  "Error on scrape endpoint", "Endpoint", endpoint.URL)
				if endpoint.Status == "Warning" {
					// check how long has it been since the last successful scrape
					// if it is more than the timeout, then set to error and stop scraping
					//endpoint.Status = "Error"
					dur, _ := time.ParseDuration(endpoint.Config.ScrapeTimeout)
					if time.Since(endpoint.LastUpdateTime) > dur {
						endpoint.Status = "Error"
						endpoint.Message = "Fail to scrape endpoint."
						endpoint.Stop()
					}
				} else {
					endpoint.Status = "Warning"
					endpoint.Message = "Fail to scrape endpoint."
				}
			}
		}
		select {
		case <-c:
			continue
		case <-endpoint.ctx.Done():
			return
		}
	}
}

func (endpoint *ExporterEndpoint) Stop() {
	// Send signal to stop
	endpoint.Status = "Stopped"
	defer endpoint.cancel()
	// Then wait
}

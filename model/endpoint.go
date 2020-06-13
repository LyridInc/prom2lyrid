package model

import (
	"context"
	"github.com/google/uuid"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"log"
	"net/http"
	"sync"
	"time"
)

// Endpoint Status Enums:
// Starting
// Started
// Warning
// Error
// Stopping
// Stopped
//

type ExporterEndpoint struct {
	ID     string       `json:"id"`
	URL    string       `json:"url"`
	Config ScrapeConfig `json:"config"`

	Status           string            `json:"status"`
	LastScrape       time.Time         `json:"last_scrape"`
	AdditionalLabels map[string]string `json:"additional_labels"`

	Message   string `json:"message"`
	IsUpdated bool   `json:"is_updated"`

	//IsScraping bool `json:-`
	//Stopping

	DurationSinceLastUpdate time.Duration
	LastUpdateTime          time.Time

	mux    sync.Mutex
	Result []*dto.MetricFamily `json:"-"`
}

func CreateEndpoint(url string) ExporterEndpoint {
	return ExporterEndpoint{
		ID:     uuid.New().String(),
		URL:    url,
		Config: CreateDefaultScrapeConfig(),
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

	duration, _ := time.ParseDuration(endpoint.Config.ScrapeInterval)
	for c := time.Tick(duration); ; {

		if endpoint.Status == "Error" {
			// do not scrape
		} else {

			log.Println("Running endpoint: " + endpoint.URL)
			start := time.Now()
			result, err := endpoint.Scrape()
			endpoint.Result = result
			log.Println("Endpoint ", endpoint.URL, " took (ms): ", time.Now().Sub(start).Milliseconds())
			if err == nil {
				endpoint.SetUpdate(true)
			} else {
				if endpoint.Status == "Warning" {
					// check how long has it been since the last successful scrape

					// if it is more than the timeout, then set to error and stop scraping
					//endpoint.Status = "Error"
				} else {
					endpoint.Status = "Warning"
				}
			}
		}
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (endpoint *ExporterEndpoint) Stop() {
	// Send signal to stop

	// Then wait
}

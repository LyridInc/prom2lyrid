package model

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"log"
	"net/http"
	"time"
)

type ExporterEndpoint struct {
	ID     string       `json:"id" binding:"required"`
	URL    string       `json:"url" binding:"required"`
	Config ScrapeConfig `json:"config" binding:"required"`

	Status           string            `json:"status" binding:"required"`
	LastScrape       time.Time         `json:"last_scrape"`
	AdditionalLabels map[string]string `json:"additional_labels"`

	Message string `json:"message"`
	//Result []*prom2json.Family `json:"result"`

}

func (endpoint *ExporterEndpoint) Scrape() []*dto.MetricFamily {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	transport.ResponseHeaderTimeout = time.Minute

	mfChan := make(chan *dto.MetricFamily, 1024)

	err := prom2json.FetchMetricFamilies(endpoint.URL, mfChan, transport)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	result := []*dto.MetricFamily{}
	for mf := range mfChan {
		result = append(result, mf)
	}

	endpoint.LastScrape = time.Now().UTC()

	return result
}

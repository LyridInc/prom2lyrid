package model

import "os"

type ScrapeConfig struct {
	ScrapeInterval string `json:"scrape_interval"`
	ScrapeTimeout  string `json:"scrape_timeout"`
}

func CreateDefaultScrapeConfig() ScrapeConfig {
	return ScrapeConfig{
		ScrapeInterval: os.Getenv("SCRAPE_INTERVAL"),
		ScrapeTimeout:  os.Getenv("SCRAPE_TIMEOUT"),
	}
}

type Credential struct {
	Key 	  	string `json:"key"`
	Secret		string `json:"secret"`
}
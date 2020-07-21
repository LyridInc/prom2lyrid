package model

type ScrapeConfig struct {
	ScrapeInterval string `json:"scrape_interval"`
	ScrapeTimeout  string `json:"scrape_timeout"`
}

func CreateDefaultScrapeConfig() ScrapeConfig {
	return ScrapeConfig{
		ScrapeInterval: "1m",
		ScrapeTimeout:  "10m",
	}
}

type Credential struct {
	Key 	  	string `json:"key"`
	Secret		string `json:"secret"`
}
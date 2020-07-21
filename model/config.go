package model

import "github.com/tkanos/gonfig"

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

func GetCredential() (Credential, error) {
	configuration := Credential{}
	err := gonfig.GetConf("credential.json", &configuration)
	return configuration, err
}
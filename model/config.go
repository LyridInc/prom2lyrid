package model

type ScrapeConfig struct {
	ScrapeInterval string `json:"scrape_interval" binding:"required"`
	ScrapeTimeout  string `json:"scrape_timeout" binding:"required"`
}

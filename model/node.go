package model

type Node struct {
	ID        string                       `json:"id" binding:"required"`
	HostName  string                       `json:"hostname" binding:"required"`
	Endpoints map[string]*ExporterEndpoint `json:"endpoints" binding:"required"`
}

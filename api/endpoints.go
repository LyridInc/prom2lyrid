package api

import (
	"github.com/gin-gonic/gin"
	"prom2lyrid/manager"
	"prom2lyrid/model"
)

func GetEndpoints(c *gin.Context) {
	c.JSON(200, manager.GetInstance().Node.Endpoints)
}

func AddEndpoints(c *gin.Context) {
	var jsonreq model.ExporterEndpoint
	if err := c.ShouldBindJSON(&jsonreq); err == nil {
		endpoint := model.CreateEndpoint(jsonreq.URL)
		mgr := manager.GetInstance()
		mgr.Node.AddEndpoint(endpoint)
		c.JSON(200, endpoint)
	} else {
		c.JSON(400, err)
	}
}

func ScrapeResult(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	endpoint := mgr.Node.Endpoints[id]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}

	if mgr.ResultCache[id] == nil {
		result, _ := endpoint.Scrape()
		mgr.ResultCache[id] = result
	}

	c.JSON(200, mgr.ResultCache[id])
}

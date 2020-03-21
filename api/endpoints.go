package api

import (
	"github.com/gin-gonic/gin"
	"prom2lyrid/manager"
)

func GetEndpoints(c *gin.Context) {
	c.JSON(200, manager.GetInstance().Node.Endpoints)
}

func ScrapeResult(c *gin.Context) {
	mgr := manager.GetInstance()
	endpoint := mgr.Node.Endpoints[c.Param("id")]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}

	if mgr.ResultCache["endpoint"] == nil {
		mgr.ResultCache["endpoint"] = endpoint.Scrape()
	}

	c.JSON(200, mgr.ResultCache["endpoint"])
}

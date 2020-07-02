package api

import (
	"github.com/gin-gonic/gin"
	"prom2lyrid/manager"
	"prom2lyrid/model"
)

//
// @Summary Get list of current an endpoints
// @Description List of current an endpoints
// @Produce json
// @Success 200 {object} map[string]model.ExporterEndpoint "returns the exporter created"
// @Failure 400 {string} string "error"
// @Router /endpoints/list [get]
// @Tags endpoints
func GetEndpoints(c *gin.Context) {

	for _, endpoint := range manager.GetInstance().Node.Endpoints {
		endpoint.SetTimeDuration()
	}
	c.JSON(200, manager.GetInstance().Node.Endpoints)
}

//
// @Summary Add an endpoint
// @Description Add an endpoint
// @Produce json
// @Param request body model.ExporterEndpoint true "adding url endpoint"
// @Success 200 {object} model.ExporterEndpoint "returns the exporter created"
// @Failure 400 {string} string "error"
// @Router /endpoints/add [post]
// @Tags endpoints
func AddEndpoints(c *gin.Context) {
	var request model.ExporterEndpoint
	if err := c.ShouldBindJSON(&request); err == nil {
		endpoint := model.CreateEndpoint(request.URL)
		mgr := manager.GetInstance()
		mgr.Node.AddEndpoint(endpoint)
		mgr.WriteConfig()
		c.JSON(200, endpoint)
	} else {
		c.JSON(400, err)
	}
}

func DeleteEndpoint(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	endpoint := mgr.Node.Endpoints[id]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}
	delete(mgr.Node.Endpoints, id)
	endpoint.Stop()
}

func UpdateEndpointLabel(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	endpoint := mgr.Node.Endpoints[id]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}

}

//
// @Summary Get scrape result of current an endpoints
// @Description Get scrape result of current an endpoints
// @Produce json
// @Param id path string true "id of exporter"
// @Success 200 {object} interface{} "returns the current scrape result created"
// @Failure 400 {string} string "error"
// @Router /endpoints/scrape/{id} [get]
// @Tags endpoints
func ScrapeResult(c *gin.Context) {
	mgr := manager.GetInstance()
	id := c.Param("id")
	endpoint := mgr.Node.Endpoints[id]

	if endpoint == nil {
		c.JSON(404, "endpoint not found")
		return
	}

	c.JSON(200, mgr.Node.Endpoints[id].Result)
}

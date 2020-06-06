package api

import (
	"github.com/gin-gonic/gin"
	"prom2lyrid/manager"
)

func Reload(c *gin.Context) {

}

func DumpConfig(c *gin.Context) {

	mgr := manager.GetInstance()
	mgr.WriteConfig()

	c.JSON(200, mgr.Node)
}

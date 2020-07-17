package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"prom2lyrid/manager"
)

func Reload(c *gin.Context) {

}

func DumpConfig(c *gin.Context) {

	mgr := manager.GetInstance()
	mgr.WriteConfig()

	c.JSON(200, mgr.Node)
}

func GetCredential(c *gin.Context) {
	c.JSON(200, map[string]string{"key": os.Getenv("LYRID_KEY"), "secret": os.Getenv("LYRID_SECRET")})
}

func SetCredential(c *gin.Context) {
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err == nil {
		os.Setenv("LYRID_KEY", request["key"])
		os.Setenv("LYRID_SECRET", request["secret"])
		c.JSON(200, map[string]string{"key": os.Getenv("LYRID_KEY"), "secret": os.Getenv("LYRID_SECRET")})
	} else {
		c.JSON(400, err)
	}
}
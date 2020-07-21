package api

import (
	"github.com/gin-gonic/gin"
	"prom2lyrid/manager"
	"prom2lyrid/model"
)

func Reload(c *gin.Context) {

}

func DumpConfig(c *gin.Context) {

	mgr := manager.GetInstance()
	mgr.WriteConfig()

	c.JSON(200, mgr.Node)
}

func GetCredential(c *gin.Context) {
	c.JSON(200, manager.GetInstance().Node.Credential)
}

func SetCredential(c *gin.Context) {
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err == nil {
		credential := model.Credential{}
		credential.Key =  request["key"]
		credential.Secret =  request["secret"]
		mgr := manager.GetInstance()
		mgr.Node.Credential = credential
		mgr.WriteConfig()
		c.JSON(200, credential)
	} else {
		c.JSON(400, err)
	}
}

func CheckLyridConnection(c *gin.Context) {
	credential := manager.GetInstance().Node.Credential
	if (len(credential.Key) > 0 && len(credential.Secret) > 0) {
		// TODO Call SDK to check connection status ...
		c.JSON(200, map[string]string{"status":"OK"})
	} else {
		c.JSON(200, map[string]string{"status":"ERROR"})
	}
}

func SetIsLocal(c *gin.Context) {
	var request map[string]bool
	if err := c.ShouldBindJSON(&request); err == nil {
		isLocal, _ := request["is_local"]
		manager.GetInstance().Node.IsLocal = isLocal
		manager.GetInstance().WriteConfig()
		c.JSON(200, manager.GetInstance().Node.IsLocal)
	} else {
		c.JSON(400, err)
	}
}

func GetIsLocal(c *gin.Context) {
	c.JSON(200, manager.GetInstance().Node.IsLocal)
}

func GetServerlessUrl(c *gin.Context) {
	c.JSON(200, manager.GetInstance().Node.ServerlessUrl)
}

func SetServerlessUrl(c *gin.Context) {
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err == nil {
		manager.GetInstance().Node.ServerlessUrl = request["url"]
		manager.GetInstance().WriteConfig()
		c.JSON(200, manager.GetInstance().Node.ServerlessUrl)
	} else {
		c.JSON(400, err)
	}
}
package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	credential, _ := model.GetCredential()
	c.JSON(200, credential)
}

func SetCredential(c *gin.Context) {
	var request map[string]string
	if err := c.ShouldBindJSON(&request); err == nil {
		credential := model.Credential{}
		credential.Key =  request["key"]
		credential.Secret =  request["secret"]
		f, _ := json.MarshalIndent(credential, "", " ")
		_ = ioutil.WriteFile("credential.json", f, 0644)
		c.JSON(200, credential)
	} else {
		c.JSON(400, err)
	}
}

func CheckLyridConnection(c *gin.Context) {
	credential, _ := model.GetCredential()
	if (len(credential.Key) > 0 && len(credential.Secret) > 0) {
		// TODO Call SDK to check connection status ...
		c.JSON(200, map[string]string{"status":"OK"})
	} else {
		c.JSON(200, map[string]string{"status":"ERROR"})
	}
}
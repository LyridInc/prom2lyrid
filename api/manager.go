package api

import (
	"github.com/LyridInc/go-sdk"
	"github.com/gin-gonic/gin"
	"os"
	"prom2lyrid/manager"
	"prom2lyrid/model"
	"prom2lyrid/utils"
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
		credential.Key = request["key"]
		credential.Secret = request["secret"]
		err = sdk.GetInstance().Initialize(credential.Key, credential.Secret)
		if err == nil {
			mgr := manager.GetInstance()
			mgr.Node.Credential = credential
			mgr.WriteConfig()
			sdk.GetInstance().ExecuteFunction(os.Getenv("FUNCTION_ID"), "LYR", utils.JsonEncode(model.LyFnInputParams{Command: "AddGateway", Gateway: mgr.Node}))
			c.JSON(200, credential)
		} else {
			c.JSON(400, err)
		}
	} else {
		c.JSON(400, err)
	}
}

func CheckLyridConnection(c *gin.Context) {
	credential := manager.GetInstance().Node.Credential
	if len(credential.Key) > 0 && len(credential.Secret) > 0 {
		// TODO Call SDK to check connection status ...
		user := sdk.GetInstance().GetUserProfile()
		if user != nil {
			account := sdk.GetInstance().GetAccountProfile()
			c.JSON(200, account)
		} else {
			c.JSON(200, map[string]string{"status": "OK"})
		}
	} else {
		c.JSON(200, map[string]string{"status": "ERROR"})
	}
}

func SetIsLocal(c *gin.Context) {
	var request map[string]bool
	if err := c.ShouldBindJSON(&request); err == nil {
		isLocal, _ := request["is_local"]
		manager.GetInstance().Node.IsLocal = isLocal
		if manager.GetInstance().Node.IsLocal {
			sdk.GetInstance().SimulateServerless(manager.GetInstance().Node.ServerlessUrl)
		} else {
			sdk.GetInstance().DisableSimulate()
		}
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
		if manager.GetInstance().Node.IsLocal {
			sdk.GetInstance().SimulateServerless(manager.GetInstance().Node.ServerlessUrl)
		}
		manager.GetInstance().WriteConfig()
		c.JSON(200, manager.GetInstance().Node.ServerlessUrl)
	} else {
		c.JSON(400, err)
	}
}

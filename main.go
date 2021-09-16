package main

import (
	"context"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"prom2lyrid/api"
	"prom2lyrid/logger"
	"prom2lyrid/manager"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API to Connect Prometheus Metrics to Lyrid Endpoint
// @version 0.0.1
// @description This is the initial definition to use Lyrid REST API
// @termsOfService https://lyrid.io/terms-of-use

// @contact.name Lyrid Support
// @contact.url https://lyrid.io
// @contact.email support@lyrid.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.GetInstance().Init()
	manager.GetInstance().Init()
	go manager.GetInstance().Run(context.Background())

	level.Info(logger.GetInstance().Logger).Log("Message", "Starting prom2Lyrid")

	router := gin.Default()
	router.Use(ginprom.PromMiddleware(nil))
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AddAllowHeaders("*")
	router.Use(cors.New(config))

	router.Use(gin.LoggerWithWriter(logger.GetInstance().LogWriter))

	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	//router.GET("/", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{"message": "home"})
	//})

	router.Use(static.Serve("/", static.LocalFile("./web/build", true)))

	manager := router.Group("/manager")
	{
		manager.POST("/reload", api.Reload)
		//manager.POST("/setup", api.Reload)
		manager.GET("/dump", api.DumpConfig)
	}

	endpoints := router.Group("/endpoints")
	{
		endpoints.GET("/list", api.GetEndpoints)
		endpoints.POST("/add", api.AddEndpoints)
		endpoints.POST("/update/:id/labels", api.UpdateEndpointLabel)
		endpoints.DELETE("/delete/:id", api.DeleteEndpoint)
		endpoints.GET("/stop/:id", api.StopEndpoint)
		endpoints.GET("/restart/:id", api.RestartEndpoint)
		//endpoints.POST("/get/:id", api.AddEndpoints)
		endpoints.GET("/scrape/:id", api.ScrapeResult)
	}
	configuration := router.Group("/config")
	{
		configuration.GET("/credential", api.GetCredential)
		configuration.POST("/credential", api.SetCredential)
		configuration.GET("/credential/status", api.CheckLyridConnection)
		configuration.GET("/serverless", api.GetServerlessUrl)
		configuration.POST("/serverless", api.SetServerlessUrl)
		configuration.GET("/local", api.GetIsLocal)
		configuration.POST("/local", api.SetIsLocal)
	}
	router.Use(static.Serve("/docs", static.LocalFile("./docs", true)))
	url := ginSwagger.URL(os.Getenv("SWAGGER_ROOT_URL") + "/docs/swagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(":" + os.Getenv("LISTENING_PORT"))
}

package main

import (
	"context"
	"fmt"
	"github.com/LyridInc/go-sdk"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"prom2lyrid/api"
	"prom2lyrid/manager"
)

func main() {
	godotenv.Load()
	fmt.Println(go_sdk.Hello())

	manager.GetInstance().Init()
	go manager.GetInstance().Run(context.Background())

	router := gin.Default()
	router.Use(ginprom.PromMiddleware(nil))
	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "home"})
	})

	manager := router.Group("/manager")
	{
		manager.POST("/reload", api.Reload)
		//manager.POST("/setup", api.Reload)
		manager.POST("/dump", api.DumpConfig)
	}

	endpoints := router.Group("/endpoints")
	{
		endpoints.GET("/list", api.GetEndpoints)
		//endpoints.GET("/get/:id", api.AddEndpoints)
		endpoints.POST("/add", api.AddEndpoints)
		//endpoints.PUT("/update", api.AddEndpoints)
		//endpoints.DELETE("/delete", api.AddEndpoints)
		//endpoints.POST("/get/:id", api.AddEndpoints)
		endpoints.GET("/scrape/:id", api.ScrapeResult)
	}

	router.Run(":8081")
}

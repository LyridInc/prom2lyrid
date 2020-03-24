package main

import (
	"fmt"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"prom2lyrid/api"
	"prom2lyrid/manager"
)

func main() {
	fmt.Println("Hello")

	godotenv.Load()

	manager.GetInstance().Init()
	go manager.GetInstance().Run()

	router := gin.Default()
	router.Use(ginprom.PromMiddleware(nil))
	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "home"})
	})

	manager := router.Group("/manager")
	{
		manager.POST("/reload", api.Reload)
	}

	endpoints := router.Group("/endpoints")
	{
		endpoints.GET("/list", api.GetEndpoints)
		endpoints.GET("/scrape/:id", api.ScrapeResult)
	}

	router.Run(":8081")
}

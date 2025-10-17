package main

import (
	"log"
	cache "main/Cache"
	getwt "main/GetWeatherTime/GetWT"
	loggerconfig "main/LoggerConfig"
	corsmiddleware "main/corsMiddleware"
	"main/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cache.Cache.Init(100)
	loggerconfig.Init()

	r := gin.Default()
	r.Use(metrics.GinMiddleware())
	r.Use(corsmiddleware.CorsMiddleware())
	r.GET("/weather", getwt.GetWT)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("Server starting at :8000")
	log.Fatal(r.Run(":8000"))
}

package main

import (
	"log"
	cache "main/Cache"
	getwt "main/GetWeatherTime/GetWT"
	loggerconfig "main/LoggerConfig"
	corsmiddleware "main/corsMiddleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cache.Cache.Init(100)
	loggerconfig.Init()

	r := gin.Default()
	r.Use(corsmiddleware.CorsMiddleware())
	r.GET("/weather", getwt.GetWT)

	log.Println("Server starting at :8000")
	log.Fatal(r.Run(":8000"))
}

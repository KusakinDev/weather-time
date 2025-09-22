package byid

import (
	"encoding/json"
	cache "main/Cache"
	owmstr "main/Struct/owmStr"
	"main/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func GetWeatherById(c *gin.Context, city int) {

	logger.Infoln("Find by id:", city)

	weatherFront, isFind := cache.Cache.FindINT(city)
	if isFind {
		c.JSON(http.StatusOK, weatherFront)
		logger.Println("Request for", city, "is already exist!")
		logger.Infoln("Response to front:", weatherFront)
		return
	}

	apiURL := config.WEATHER_API_3 + strconv.Itoa(city) + config.WEATHER_API_2

	logger.Println("Reguest on:", apiURL)

	response, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "OpenWeatherAPI error"})
		logger.Errorln("OpenWeatherAPI error")
		return
	}
	defer response.Body.Close()

	var weatherData owmstr.WeatherData
	err = json.NewDecoder(response.Body).Decode(&weatherData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect city name"})
		logger.Errorln("Incorrect city name")
		return
	}

	weatherFront.Id = weatherData.Syst.Id
	weatherFront.Cod = weatherData.Cod
	weatherFront.Name = weatherData.Name
	weatherFront.Timezone = weatherData.Timezone / (60 * 60)
	weatherFront.Temp = weatherData.Main.Temp
	weatherFront.Icon = config.IMAGE_API_1 + weatherData.Weather[0].Icon + config.IMAGE_API_2

	logger.Infoln("Response to front:", weatherFront)
	cache.Cache.InsertINT(weatherFront.Id, weatherFront)

	c.JSON(http.StatusOK, weatherFront)
	logger.Infoln("Response has been sent")
}

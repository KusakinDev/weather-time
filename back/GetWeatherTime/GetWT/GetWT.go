package getwt

import (
	byid "main/GetWeatherTime/byId"
	byname "main/GetWeatherTime/byName"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func swithType(city string) (int, int) {
	if num, err := strconv.Atoi(city); err == nil {
		return num, 0
	} else {
		return -1, 1
	}
}

func GetWT(c *gin.Context) {

	city := c.Query("name")
	logger.Infoln(city)

	if city == "" {
		c.JSON(400, gin.H{"message": "Name param is missing"})
		return
	}

	cityInt, check := swithType(city)
	switch check {
	case 0:
		byid.GetWeatherById(c, cityInt)
	case 1:
		byname.GetWeatherByName(c, city)
	default:

	}
}

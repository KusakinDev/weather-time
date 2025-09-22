package owmstr

type WeatherData struct {
	Main     MainInfo  `json:"main"`
	Weather  []Weather `json:"weather"`
	Timezone int       `json:"timezone"`
	Name     string    `json:"name"`
	Cod      int       `json:"cod"`
	Syst     Sys       `json:"sys"`
}

type MainInfo struct {
	Temp float32 `json:"temp"`
}

type Weather struct {
	Icon string `json:"icon"`
}

type Sys struct {
	Id int `json:"id"`
}

/*
{
    "weather": [
        {
            "icon": "01d"
        }
    ],
    "main": {
        "temp": -6.39,
    },
    "timezone": 25200,
    "name": "Новосибирск",
    "cod": 200
    "sys": {
        "id": 197864,
    },

}
*/

/*
{
    "coord": {
        "lon": 82.9344,
        "lat": 55.0411
    },
    "weather": [
        {
            "id": 800,
            "main": "Clear",
            "description": "ясно",
            "icon": "01d"
        }
    ],
    "base": "stations",
    "main": {
        "temp": -6.39,
        "feels_like": -11.18,
        "temp_min": -6.39,
        "temp_max": -6.39,
        "pressure": 1037,
        "humidity": 79,
        "sea_level": 1037,
        "grnd_level": 1015
    },
    "visibility": 10000,
    "wind": {
        "speed": 3,
        "deg": 170
    },
    "clouds": {
        "all": 0
    },
    "dt": 1739783387,
    "sys": {
        "type": 1,
        "id": 8958,
        "country": "RU",
        "sunrise": 1739756821,
        "sunset": 1739792269
    },
    "timezone": 25200,
    "id": 1496747,
    "name": "Новосибирск",
    "cod": 200
}
*/

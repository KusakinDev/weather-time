package responsestr

type WeatherToFront struct {
	Cod      int     `json:"cod"`
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Temp     float32 `json:"temp"`
	Timezone int     `json:"timezone"`
	Icon     string  `json:"icon"`
}

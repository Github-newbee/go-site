package v1

type WeatherRequest struct {
	City string `form:"city"`
}

type WeatherResponse struct {
	Response
	Data interface{}
}

// 高德地图天气API响应结构
type AmapWeatherResponse struct {
	Status   string            `json:"status"`
	Count    string            `json:"count"`
	Info     string            `json:"info"`
	InfoCode string            `json:"infocode"`
	Lives    []AmapWeatherLive `json:"lives"`
}

type AmapWeatherLive struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	AdCode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	WindDirection    string `json:"winddirection"`
	WindPower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	ReportTime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	HumidityFloat    string `json:"humidity_float"`
}

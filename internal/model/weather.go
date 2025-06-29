package model

type Weather struct {
	BaseModel

	AdCode        string `json:"adcode" gorm:"unique;not null; type:varchar(100); default:''; comment:城市编码"`
	Weather       string `json:"weather" gorm:"type:varchar(100); default:''; comment:天气"`
	Temperature   string `json:"temperature" gorm:"type:varchar(100); default:''; comment:温度"`
	WindPower     string `json:"windpower" gorm:"type:varchar(100); default:''; comment:风力"`
	Province      string `json:"province" gorm:"type:varchar(100); default:''; comment:省份"`
	ReportTime    string `json:"reporttime" gorm:"type:varchar(100); default:''; comment:报告时间"`
	Humidity      string `json:"humidity" gorm:"type:varchar(100); default:''; comment:湿度"`
	City          string `json:"city" gorm:"type:varchar(100); default:''; comment:城市"`
	WindDirection string `json:"winddirection" gorm:"type:varchar(100); default:''; comment:风向"`
}

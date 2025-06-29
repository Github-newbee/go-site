package system

import (
	"go-site/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterWeatherRoutes(group *gin.RouterGroup, handlers *handler.Handler) {
	v1 := group.Group("/weather")
	{
		v1.GET("", handlers.WeatherHandler.GetWeather) // 获取天气数据
	}
}

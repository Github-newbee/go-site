package handler

import (
	v1 "go-site/api/v1"
	"go-site/internal/service"
	"go-site/pkg/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weatherService service.WeatherService
}

func NewWeatherHandler(handler *Handler, weatherService service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(ctx *gin.Context) {
	req := v1.WeatherRequest{}
	queryErr := request.Assign(ctx, &req)
	if queryErr != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, queryErr.Error())
		return
	}

	weather, err := h.weatherService.GetWeather(req)

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, weather)
}

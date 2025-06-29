package handler

import (
	"go-site/internal/service"
	"go-site/internal/service/common"
	"go-site/pkg/jwt"
	"go-site/pkg/log"
	"go-site/pkg/sid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger          *log.Logger
	UserHandler     *UserHandler
	CategoryHandler *CategoryHandler
	WebSiteHandler  *WebsiteHandler
	FileHandler     *FileHandler
	WeatherHandler  *WeatherHandler
}

func NewHandler(
	logger *log.Logger,
	userService service.UserService,
	categoryService service.CategoryService,
	webSiteService service.WebsiteService,
	fileService common.FileService,
	weatherService service.WeatherService,
) *Handler {
	h := &Handler{}
	// 移除对 Handler 的依赖
	h.UserHandler = NewUserHandler(h, userService)
	h.CategoryHandler = NewCategoryHandler(h, categoryService)
	h.WebSiteHandler = NewWebsiteHandler(h, webSiteService)
	h.FileHandler = NewFileHandler(h, fileService)
	h.WeatherHandler = NewWeatherHandler(h, weatherService)
	return h
}

func GetUserIdFromCtx(ctx *gin.Context) sid.SnowflakeID {
	v, exists := ctx.Get("claims")
	if !exists {
		return sid.SnowflakeID(0)
	}
	id, err := sid.NewSnowflakeIDFromString(v.(*jwt.MyCustomClaims).UserId)
	if err != nil {
		return sid.SnowflakeID(0)
	}
	return id
}

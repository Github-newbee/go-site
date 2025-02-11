package handler

import (
	"go-my-demo/internal/service"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger          *log.Logger
	UserHandler     *UserHandler
	CategoryHandler *CategoryHandler
	WebSiteHandler  *WebsiteHandler
}

func NewHandler(
	logger *log.Logger,
	userService service.UserService,
	categoryService service.CategoryService,
	webSiteService service.WebsiteService,
) *Handler {
	h := &Handler{}
	// 移除对 Handler 的依赖
	h.UserHandler = NewUserHandler(h, userService)
	h.CategoryHandler = NewCategoryHandler(h, categoryService)
	h.WebSiteHandler = NewWebsiteHandler(h, webSiteService)
	return h
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}

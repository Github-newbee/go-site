package handler

import (
	"go-my-demo/internal/service"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger      *log.Logger
	UserHandler *UserHandler
}

func NewHandler(
	logger *log.Logger,
	userService service.UserService,
) *Handler {
	h := &Handler{}
	h.UserHandler = NewUserHandler(h, userService) // 移除对 Handler 的依赖
	return h
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}

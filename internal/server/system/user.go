package system

import (
	"go-my-demo/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(group *gin.RouterGroup, handlers *handler.Handler) {
	v1 := group.Group("/user")
	{
		v1.GET("", handlers.UserHandler.GetAllUsers)        // 获取列表
		v1.GET("/profile", handlers.UserHandler.GetProfile) // 获取个人信息
		v1.GET("/:id", handlers.UserHandler.GetProfileByID)
		v1.PUT("", handlers.UserHandler.UpdateProfile) // 更新用户
	}
}

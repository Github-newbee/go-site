package system

import (
	"go-my-demo/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterWebsiteRoutes(group *gin.RouterGroup, handlers *handler.Handler) {
	v1 := group.Group("/website")
	{
		// v1.GET("", handlers.WebSiteHandler) // 获取列表
		v1.POST("", handlers.WebSiteHandler.CreateWebsite) // 创建
		// v1.GET("/:id", handlers.UserHandler.GetProfileByID)
		// v1.PUT("/:id", handlers.CategoryHandler.UpdateCategory)
	}
}

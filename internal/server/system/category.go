package system

import (
	"go-my-demo/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(group *gin.RouterGroup, handlers *handler.Handler) {
	v1 := group.Group("/category")
	{
		v1.GET("", handlers.CategoryHandler.GetAllCategory)  // 获取列表
		v1.POST("", handlers.CategoryHandler.CreateCategory) // 创建
		// v1.GET("/:id", handlers.UserHandler.GetProfileByID)
		v1.PUT("/:id", handlers.CategoryHandler.UpdateCategory)
	}
}

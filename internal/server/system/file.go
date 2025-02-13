package system

import (
	"go-my-demo/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(group *gin.RouterGroup, handlers *handler.Handler) {
	v1 := group.Group("/file")
	{
		v1.POST("/upload", handlers.FileHandler.UploadFile)
	}
}

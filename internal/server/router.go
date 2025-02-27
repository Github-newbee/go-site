package server

import (
	"go-my-demo/internal/handler"
	"go-my-demo/internal/middleware"
	"go-my-demo/internal/server/system"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Router struct {
	handlers *handler.Handler
	jwt      *jwt.JWT
	logger   *log.Logger
	conf     *viper.Viper
}

// ProvideRouter wire provider
func ProvideRouter(
	handlers *handler.Handler,
	jwt *jwt.JWT,
	logger *log.Logger,
	conf *viper.Viper,
) *Router {
	return &Router{
		handlers: handlers,
		jwt:      jwt,
		logger:   logger,
		conf:     conf,
	}
}

func (r *Router) Register(e *gin.Engine) {

	v1 := e.Group("/v1")
	{
		v1.POST("login", r.handlers.UserHandler.Login)
		v1.POST("register", r.handlers.UserHandler.Register)

		r.registerSystemRoutes(v1)

		r.registerFileRoutes(v1)
	}

}

// 系统路由
func (r *Router) registerSystemRoutes(group *gin.RouterGroup) {
	systemGroup := group.Group("/system")
	// 严格认证
	systemGroup.Use(middleware.StrictAuth(r.jwt, r.logger))
	{
		// 用户路由
		system.RegisterUserRoutes(systemGroup, r.handlers)
		// 分类路由
		system.RegisterCategoryRoutes(systemGroup, r.handlers)
		// 网站路由
		system.RegisterWebsiteRoutes(systemGroup, r.handlers)
	}
}

// 文件路由
func (r *Router) registerFileRoutes(group *gin.RouterGroup) {
	fileGroup := group.Group("/")
	// 严格认证
	fileGroup.Use(middleware.StrictAuth(r.jwt, r.logger))
	{
		system.RegisterFileRoutes(group, r.handlers)
	}
}

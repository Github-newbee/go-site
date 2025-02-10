package server

import (
	apiV1 "go-my-demo/api/v1"
	"go-my-demo/internal/middleware"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"
	"go-my-demo/pkg/server/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	router *Router,
) *http.Server {
	// DebugMode ReleaseMode TestMode
	gin.SetMode(gin.ReleaseMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Hello, World!",
		})
	})

	s.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogMiddleware(logger),
		middleware.ResponseLogMiddleware(logger),
	)

	router.Register(s.Engine)

	return s
}

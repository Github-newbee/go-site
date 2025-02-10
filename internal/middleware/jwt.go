package middleware

import (
	v1 "go-my-demo/api/v1"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Auth 认证中间件

// StrictAuth 严格认证
func StrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求头中的 Authorization 字段
		tokenString := ctx.Request.Header.Get("Authorization")
		// 从请求头中获取 Authorization 字段的值。如果没有找到令牌，记录警告日志并返回未授权错误。
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			// 终止请求
			ctx.Abort()
			return
		}

		// 将解析出的声明（claims）存储在上下文中，并调用 recoveryLoggerFunc 记录日志。

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		// 然后调用 ctx.Next() 继续处理请求。
		ctx.Next()
	}
}

// NoStrictAuth 非严格认证
func NoStrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 尝试从请求头中的 Authorization 字段获取令牌
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			// 如果请求头中没有令牌，尝试从 Cookie 中获取令牌
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			// 如果 Cookie 中没有令牌，尝试从查询参数中获取令牌
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			// 如果没有找到令牌，继续处理请求
			ctx.Next()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.String("UserId", userInfo.UserId))
	}
}

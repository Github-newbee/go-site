package service

import (
	"context"
	"go-site/internal/repository"
	"go-site/pkg/jwt"
	"go-site/pkg/log"
	"go-site/pkg/sid"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
	rdb    *redis.Client   // 添加 Redis 客户端
	ctx    context.Context // 添加上下文
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	rdb *redis.Client, // 添加 Redis 参数
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
		rdb:    rdb,
		ctx:    context.Background(),
	}
}

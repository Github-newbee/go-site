package job

import (
	"go-my-demo/internal/repository"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"
	"go-my-demo/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}

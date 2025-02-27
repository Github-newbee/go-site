// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-my-demo/internal/handler"
	"go-my-demo/internal/job"
	"go-my-demo/internal/repository"
	"go-my-demo/internal/server"
	"go-my-demo/internal/service"
	"go-my-demo/internal/service/common"
	"go-my-demo/pkg/app"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"
	"go-my-demo/pkg/server/http"
	"go-my-demo/pkg/sid"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	jwtJWT := jwt.NewJwt(viperViper)
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	categoryRepository := repository.NewCategoryRepository(repositoryRepository)
	categoryService := service.NewCategoryService(serviceService, categoryRepository)
	websiteRepository := repository.NewWebsiteRepository(repositoryRepository)
	websiteService := service.NewWebsiteService(serviceService, websiteRepository, categoryRepository)
	fileService := common.NewFileService(serviceService)
	handlerHandler := handler.NewHandler(logger, userService, categoryService, websiteService, fileService)
	router := server.ProvideRouter(handlerHandler, jwtJWT, logger, viperViper)
	httpServer := server.NewHTTPServer(logger, viperViper, jwtJWT, router)
	jobJob := job.NewJob(transaction, logger, sidSid)
	userJob := job.NewUserJob(jobJob, userRepository)
	jobServer := server.NewJobServer(logger, userJob)
	appApp := newApp(httpServer, jobServer, viperViper)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository, repository.NewCategoryRepository, repository.NewWebsiteRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService, service.NewCategoryService, service.NewWebsiteService, common.NewFileService)

var handlerSet = wire.NewSet(handler.NewHandler)

var jobSet = wire.NewSet(job.NewJob, job.NewUserJob)

var serverSet = wire.NewSet(server.NewHTTPServer, server.NewJobServer, server.ProvideRouter)

// build App
// httpServer 和 jobServer 作为参数传入，来源于 serverSet
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,

	conf *viper.Viper,
) *app.App {
	return app.NewApp(app.WithServer(httpServer, jobServer), app.WithName(conf.GetString("app.name")))
}

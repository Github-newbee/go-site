//go:build wireinject
// +build wireinject

package wire

import (
	"go-site/internal/handler"
	"go-site/internal/job"
	"go-site/internal/repository"
	"go-site/internal/server"
	"go-site/internal/service"
	"go-site/internal/service/common"
	"go-site/pkg/app"
	"go-site/pkg/jwt"
	"go-site/pkg/log"
	"go-site/pkg/server/http"
	"go-site/pkg/sid"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// 定义依赖注入规则

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewCategoryRepository,
	repository.NewWebsiteRepository,
	repository.NewWeatherRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewCategoryService,
	service.NewWebsiteService,
	service.NewWeatherService,
	common.NewFileService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
	server.ProvideRouter,
)

// build App
// httpServer 和 jobServer 作为参数传入，来源于 serverSet
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
	conf *viper.Viper,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName(conf.GetString("app.name")),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}

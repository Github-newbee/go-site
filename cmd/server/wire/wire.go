//go:build wireinject
// +build wireinject

package wire

import (
	"go-my-demo/internal/handler"
	"go-my-demo/internal/job"
	"go-my-demo/internal/repository"
	"go-my-demo/internal/server"
	"go-my-demo/internal/service"
	"go-my-demo/pkg/app"
	"go-my-demo/pkg/jwt"
	"go-my-demo/pkg/log"
	"go-my-demo/pkg/server/http"
	"go-my-demo/pkg/sid"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// 定义依赖注入规则

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewCategoryRepository,
	repository.NewWebsiteRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewCategoryService,
	service.NewWebsiteService,
	service.NewFileService,
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

package main

import (
	"context"
	"flag"
	"fmt"

	"go-my-demo/cmd/server/wire"
	"go-my-demo/internal/model"
	"go-my-demo/internal/repository"
	"go-my-demo/pkg/config"
	"go-my-demo/pkg/log"

	"go.uber.org/zap"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// 启动环境指定，默认../../config/dev.yml 可以通过 -conf xxx路径名  指定
	var envConf = flag.String("conf", "./config/dev.yml", "config path, eg: -conf ./config/dev.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	// Initialize database connection
	db := repository.NewDB(conf, logger)
	// Auto migrate models
	err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Website{},
	)
	if err != nil {
		logger.Fatal("Auto migration failed", zap.Error(err))
	}
	app, cleanup, err := wire.NewWire(conf, logger)

	defer cleanup()
	if err != nil {
		panic(err)
	}
	fmt.Println("当前环境：", conf.GetString("env"))
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	"go-my-demo/internal/repository"
	"go-my-demo/pkg/config"
	"go-my-demo/pkg/log"
)

func main() {
	var (
		envConf = flag.String("conf", "./config/dev.yml", "config path, eg: -conf ./config/dev.yml")
		action  = flag.String("action", "up", "migration action: up, down, reset, status")
		help    = flag.Bool("help", false, "show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 加载配置
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	// 初始化数据库连接
	db := repository.NewDB(conf, logger)

	// 创建迁移器
	migrator := NewMigrator(db, logger)

	// 执行迁移操作
	switch *action {
	case "up":
		if err := migrator.Up(); err != nil {
			logger.Error(fmt.Sprintf("Migration up failed: %v", err))
			os.Exit(1)
		}
		fmt.Println("✅ Migration completed successfully!")
	case "down":
		if err := migrator.Down(); err != nil {
			logger.Error(fmt.Sprintf("Migration down failed: %v", err))
			os.Exit(1)
		}
		fmt.Println("✅ Migration rollback completed!")
	case "reset":
		if err := migrator.Reset(); err != nil {
			logger.Error(fmt.Sprintf("Migration reset failed: %v", err))
			os.Exit(1)
		}
		fmt.Println("✅ Database reset completed!")
	case "status":
		migrator.Status()
	default:
		fmt.Printf("❌ Unknown action: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`
	Database Migration Tool
	
	Usage:
		go run cmd/migrate/main.go [options]

	Options:
		-conf string    Configuration file path (default "./config/dev.yml")
		-action string  Migration action (default "up")
		-help          Show this help message

	Actions:
		up      Run migrations (create/update tables)
		down    Rollback migrations (drop tables)
		reset   Reset database (drop and recreate all tables)
		status  Show migration status

	Examples:
		go run cmd/migrate/main.go -action=up
		go run cmd/migrate/main.go -action=status -conf=./config/prod.yml
		go run cmd/migrate/main.go -action=reset
	`)
}

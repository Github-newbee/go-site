package migrator

import (
	"fmt"
	"go-site/internal/model"
	"go-site/pkg/log"

	"gorm.io/gorm"
)

type Migrator struct {
	db     *gorm.DB
	logger *log.Logger
	models []interface{}
}

func NewMigrator(db *gorm.DB, logger *log.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
		models: []interface{}{
			&model.User{},
			&model.Category{},
			&model.Website{},
			&model.Weather{},
		},
	}
}

// Up 执行数据库迁移
func (m *Migrator) Up() error {
	m.logger.Info("Starting database migration...")

	for _, model := range m.models {
		modelName := fmt.Sprintf("%T", model)
		m.logger.Info("Migrating model")

		if err := m.db.AutoMigrate(model); err != nil {
			m.logger.Error("Failed to migrate model")
			return err
		}

		fmt.Printf("✅ Migrated: %s\n", modelName)
	}

	m.logger.Info("Database migration completed")
	return nil
}

// Down 回滚数据库迁移（删除表）
func (m *Migrator) Down() error {
	m.logger.Info("Starting database rollback...")

	// 反向删除表（避免外键约束问题）
	for i := len(m.models) - 1; i >= 0; i-- {
		model := m.models[i]
		modelName := fmt.Sprintf("%T", model)

		m.logger.Info("Dropping table for model")

		if err := m.db.Migrator().DropTable(model); err != nil {
			m.logger.Error("Failed to drop table")
			return err
		}

		fmt.Printf("✅ Dropped table: %s\n", modelName)
	}

	m.logger.Info("Database rollback completed")
	return nil
}

// Reset 重置数据库（删除后重新创建）
func (m *Migrator) Reset() error {
	fmt.Println("⚠️  This will delete all data! Are you sure? (y/N)")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("Operation cancelled")
		return nil
	}

	m.logger.Info("Resetting database...")

	// 先删除所有表
	if err := m.Down(); err != nil {
		return err
	}

	// 再重新创建
	if err := m.Up(); err != nil {
		return err
	}

	m.logger.Info("Database reset completed")
	return nil
}

// Status 显示迁移状态
func (m *Migrator) Status() {
	fmt.Println("\n📊 Database Migration Status:")
	fmt.Println("================================")

	for _, model := range m.models {
		modelName := fmt.Sprintf("%T", model)

		if m.db.Migrator().HasTable(model) {
			fmt.Printf("✅ %s -> Table exists\n", modelName)

			// 检查列信息
			if columns, err := m.db.Migrator().ColumnTypes(model); err == nil {
				fmt.Printf("   Columns: %d\n", len(columns))
			}
		} else {
			fmt.Printf("❌ %s -> Table missing\n", modelName)
		}
	}

	fmt.Println("================================")
}

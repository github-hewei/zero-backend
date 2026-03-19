package runner

import (
	"context"
	"path/filepath"

	"zero-backend/pkg/logger"
	"zero-backend/pkg/migrate"

	"gorm.io/gorm"
)

// MigrateRunner 迁移执行器
type MigrateRunner struct {
	logger logger.Logger
	db     *gorm.DB
}

// NewMigrateRunner 创建迁移执行器
func NewMigrateRunner(l logger.Logger, db *gorm.DB) *MigrateRunner {
	return &MigrateRunner{
		logger: l,
		db:     db,
	}
}

// Up 执行数据库迁移
func (r *MigrateRunner) Up(filePath string) error {
	// 确定 SQL 文件路径
	sqlFile := filePath
	if sqlFile == "" {
		sqlFile = filepath.Join("data", "database.sql")
	}

	// 创建迁移器
	migrator := migrate.NewMigrator(
		r.db,
		sqlFile,
		migrate.WithLogger(r.logger),
	)

	// 执行迁移
	ctx := context.Background()
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}

	r.logger.Info("数据库迁移完成")
	return nil
}

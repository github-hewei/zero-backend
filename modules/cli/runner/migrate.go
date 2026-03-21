package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"zero-backend/internal/apperror"
	"zero-backend/internal/logger"
	"zero-backend/pkg/migrate"

	"gorm.io/gorm"
)

// MigrateRunner 迁移执行器
type MigrateRunner struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewMigrateRunner 创建迁移执行器
func NewMigrateRunner(l *logger.Logger, db *gorm.DB) *MigrateRunner {
	return &MigrateRunner{
		logger: l,
		db:     db,
	}
}

// Up 执行数据库迁移
func (r *MigrateRunner) Up(ctx context.Context, filePath string) error {
	// 确定 SQL 文件路径
	sqlFile := filePath
	if sqlFile == "" {
		sqlFile = filepath.Join("data", "database.sql")
	}

	// 创建迁移器
	migrator := migrate.NewMigrator(r.db, sqlFile, migrate.WithLogger(r.logger))

	// 执行迁移
	if err := migrator.Migrate(ctx); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return apperror.NewUserError(fmt.Sprintf("找不到指定 SQL 文件：%s", sqlFile))
		}

		var migrationError *migrate.MigrationError
		if errors.As(err, &migrationError) {
			return apperror.NewUserError(fmt.Sprintf("执行迁移失败：%s", migrationError.Error()))
		}

		return apperror.NewSystemError(err, "执行迁移失败")
	}

	r.logger.Info("数据库迁移完成")
	return nil
}

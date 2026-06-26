package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/migrate"
	"github.com/241x/zero-web/errcode"
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
			return apperror.New(errcode.NotFound, apperror.WithMsg(fmt.Sprintf("找不到指定 SQL 文件：%s", sqlFile)))
		}

		if migrationError, ok := errors.AsType[*migrate.MigrationError](err); ok {
			return apperror.New(errcode.Internal, apperror.WithMsg(fmt.Sprintf("执行迁移失败：%s", migrationError.Error())))
		}

		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("执行迁移失败"))
	}

	r.logger.Info("数据库迁移完成")
	return nil
}

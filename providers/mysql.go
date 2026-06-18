package providers

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mysql"

	"github.com/google/wire"
	"gorm.io/gorm/logger"
)

// NewMySQLConfig 从全局配置中提取 MySQL 配置
func NewMySQLConfig(cfg *config.Config) mysql.Config {
	return mysql.Config{
		Dsn:    cfg.MySQL.Dsn,
		Prefix: cfg.MySQL.Prefix,
	}
}

// MySQLProviderSet 提供MySQL数据库依赖集合
var MySQLProviderSet = wire.NewSet(
	NewMySQLConfig,
	mysql.NewDB,
	gormutil.NewLogger,
	wire.Bind(new(logger.Interface), new(*gormutil.Logger)),
)

package providers

import (
	"zero-backend/internal/storage/mysql"

	"github.com/google/wire"
	"gorm.io/gorm/logger"
)

// MySQLProviderSet 提供MySQL数据库依赖集合
var MySQLProviderSet = wire.NewSet(
	mysql.NewDB,
	mysql.NewLogger,
	wire.Bind(new(logger.Interface), new(*mysql.Logger)),
)

package providers

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/mongodb"

	"github.com/google/wire"
)

// NewMongoDBConfig 从全局配置中提取 MongoDB 配置
func NewMongoDBConfig(cfg *config.Config) mongodb.Config {
	return mongodb.Config{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
		Enabled:  cfg.MongoDB.Enabled,
	}
}

// MongoDBProviderSet 提供MongoDB数据库依赖集合
var MongoDBProviderSet = wire.NewSet(
	NewMongoDBConfig,
	mongodb.NewConn,
	wire.FieldsOf(new(*mongodb.Conn), "Client", "DB"),
)

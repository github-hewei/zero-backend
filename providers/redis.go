package providers

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/redis"
	"github.com/google/wire"
)

// NewRedisConfig 从全局配置中提取 Redis 配置
func NewRedisConfig(cfg *config.Config) redis.Config {
	return redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
}

// RedisProviderSet 提供Redis依赖集合
var RedisProviderSet = wire.NewSet(
	NewRedisConfig,
	redis.New,
)

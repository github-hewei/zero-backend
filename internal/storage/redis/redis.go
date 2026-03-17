package redis

import (
	"fmt"
	"zero-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

// New 新建Redis客户端
func New(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client
}

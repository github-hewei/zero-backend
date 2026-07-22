package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/provider"
	"zero-backend/internal/worker"

	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
)

func main() {
	config.Init()

	rdb := redis.New(provider.LoadRedisConfig())
	manager := queue.NewQueueManager(rdb)

	conn := mongodb.MustNewConn(provider.LoadMongoConfig())
	log := provider.NewLogger(conn.DB, "worker.log")

	registry := worker.NewRegistry(log)
	registry.Register("example", worker.NewExampleHandler(log))

	// 启动服务
	worker.NewServer(manager, registry, log).Run()
}

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
	mgr := queue.NewQueueManager(rdb)

	conn := mongodb.MustNewConn(provider.LoadMongoConfig())
	l := provider.LoadLogger(conn.DB, "worker.log")

	registry := worker.NewRegistry(l)
	registry.Register("example", &worker.ExampleHandler{L: l})

	worker.NewServer(mgr, registry, l).Run()
}

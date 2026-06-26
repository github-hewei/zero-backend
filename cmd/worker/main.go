package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/worker"
	"zero-backend/providers"

	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
)

func main() {
	cfg := config.New()

	rdb := redis.New(providers.NewRedisConfig(cfg))
	mgr := queue.NewQueueManager(rdb)

	mongoCfg := providers.NewMongoDBConfig(cfg)
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := providers.ProvideLogger(cfg.Logger, conn.DB)

	registry := worker.NewRegistry(l)
	registry.Register("example", &worker.ExampleHandler{L: l})

	worker.NewServer(mgr, registry, l).Run()
}

package main

import (
	"zero-backend/internal/app"
	"zero-backend/internal/config"
	"zero-backend/internal/worker"

	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
)

func main() {
	cfg := config.New()

	rdb := redis.New(app.NewRedisConfig(cfg))
	mgr := queue.NewQueueManager(rdb)

	mongoCfg := app.NewMongoDBConfig(cfg)
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := app.ProvideLogger(cfg.Logger, conn.DB)

	registry := worker.NewRegistry(l)
	registry.Register("example", &worker.ExampleHandler{L: l})

	worker.NewServer(mgr, registry, l).Run()
}

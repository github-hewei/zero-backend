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
	config.Init()

	rdb := redis.New(app.LoadRedisConfig())
	mgr := queue.NewQueueManager(rdb)

	conn := app.Must(mongodb.NewConn(app.LoadMongoConfig()))
	l := app.LoadLogger(conn.DB)

	registry := worker.NewRegistry(l)
	registry.Register("example", &worker.ExampleHandler{L: l})

	worker.NewServer(mgr, registry, l).Run()
}

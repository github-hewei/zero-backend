package main

import (
	"zero-backend/internal/cli"
	"zero-backend/internal/cli/runner"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/rbac"
	"zero-backend/providers"

	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
)

func main() {
	cfg := config.New()

	mongoCfg := providers.NewMongoDBConfig(cfg)
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := providers.ProvideLogger(cfg.Logger, conn.DB)

	rdb := redis.New(providers.NewRedisConfig(cfg))

	gormLog := gormutil.NewLogger(l)
	db, err := mysql.NewDB(providers.NewMySQLConfig(cfg), gormLog)
	if err != nil {
		panic(err)
	}

	app := cli.New(l, rdb)
	app.AddCommand(cli.MigrateCmd(db, l))
	app.AddCommand(cli.QueueCmd(queue.NewQueueManager(rdb)))
	app.AddCommand(cli.SyncApiCmd(runner.NewSyncApiRunner(l, rbac.NewRbacApiRepository(db))))
	app.Run()
}

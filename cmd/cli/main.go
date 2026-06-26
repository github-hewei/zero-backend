package main

import (
	"zero-backend/internal/app"
	"zero-backend/internal/cli"
	"zero-backend/internal/cli/runner"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/rbac"

	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
)

func main() {
	cfg := config.New()

	mongoCfg := app.NewMongoDBConfig(cfg)
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := app.ProvideLogger(cfg.Logger, conn.DB)

	rdb := redis.New(app.NewRedisConfig(cfg))

	gormLog := gormutil.NewLogger(l)
	db, err := mysql.NewDB(app.NewMySQLConfig(cfg), gormLog)
	if err != nil {
		panic(err)
	}

	cliApp := cli.New(l, rdb)
	cliApp.AddCommand(cli.MigrateCmd(db, l))
	cliApp.AddCommand(cli.QueueCmd(queue.NewQueueManager(rdb)))
	cliApp.AddCommand(cli.SyncApiCmd(runner.NewSyncApiRunner(l, rbac.NewRbacApiRepository(db))))
	cliApp.Run()
}

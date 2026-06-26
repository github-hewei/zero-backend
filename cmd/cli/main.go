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
	config.Init()

	conn := app.Must(mongodb.NewConn(app.LoadMongoConfig()))
	l := app.LoadLogger(conn.DB)

	rdb := redis.New(app.LoadRedisConfig())

	gormLog := gormutil.NewLogger(l)
	db := app.Must(mysql.NewDB(app.LoadMySQLConfig(), gormLog))

	cliApp := cli.New(l, rdb)
	cliApp.AddCommand(cli.MigrateCmd(db, l))
	cliApp.AddCommand(cli.QueueCmd(queue.NewQueueManager(rdb)))
	cliApp.AddCommand(cli.SyncApiCmd(runner.NewSyncApiRunner(l, rbac.NewRbacApiRepository(db))))
	cliApp.Run()
}

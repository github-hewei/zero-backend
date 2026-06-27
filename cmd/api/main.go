package main

import (
	"zero-backend/internal/api"
	"zero-backend/internal/app"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/user"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/server"
)

func main() {
	config.Init()

	conn := app.Must(mongodb.NewConn(app.LoadMongoConfig()))
	l := app.LoadLogger(conn.DB, "api.log")

	gormLog := gormutil.NewLogger(l)
	db := app.Must(mysql.NewDB(app.LoadMySQLConfig(), gormLog))

	v := bind.NewValidate()
	t := app.Must(bind.NewTrans(v))
	binder := bind.New(v, t, app.ProvideBindErrCode())

	rdb := redis.New(app.LoadRedisConfig())
	authCfg := app.Must(user.LoadConfig())

	engine := api.NewGin(l, db, binder, rdb, authCfg)

	srv := server.New(app.LoadApiServerConfig(), engine, l, app.ProvideServerOptions()...)
	srv.Run()
}

package main

import (
	"zero-backend/internal/api"
	"zero-backend/internal/app"
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/server"
)

func main() {
	config.Init()

	mongoCfg := app.LoadMongoConfig()
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := app.LoadLogger(conn.DB)

	gormLog := gormutil.NewLogger(l)
	db, err := mysql.NewDB(app.LoadMySQLConfig(), gormLog)
	if err != nil {
		panic(err)
	}

	v := bind.NewValidate()
	t, err := bind.NewTrans(v)
	if err != nil {
		panic(err)
	}
	binder := bind.New(v, t, app.ProvideBindErrCode())

	rdb := redis.New(app.LoadRedisConfig())

	srv := server.New(
		app.LoadApiServerConfig(),
		api.NewGin(l, db, binder, rdb),
		l,
		app.ProvideServerOptions()...,
	)
	srv.Run()
}

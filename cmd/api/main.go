package main

import (
	"zero-backend/internal/api"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/user"
	"zero-backend/internal/provider"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/server"
)

func main() {
	config.Init()

	conn := mongodb.MustNewConn(provider.LoadMongoConfig())
	l := provider.NewLogger(conn.DB, "api.log")

	gormLog := gormutil.NewLogger(l)
	db := mysql.MustNewDB(provider.LoadMySQLConfig(), gormLog)

	v := bind.NewValidate()
	t := bind.MustNewTrans(v)
	binder := bind.New(v, t, provider.ProvideBindErrCode())

	rdb := redis.New(provider.LoadRedisConfig())
	authCfg := user.MustLoadConfig()

	engine := api.NewGin(l, db, binder, rdb, authCfg)

	srv := server.New(provider.LoadApiServerConfig(), engine, l, provider.ProvideServerOptions()...)
	srv.Run()
}

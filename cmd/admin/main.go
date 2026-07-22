package main

import (
	"zero-backend/internal/admin"
	"zero-backend/internal/config"
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
	log := provider.LoadLogger(conn.DB, "admin.log")

	gormLog := gormutil.NewLogger(log)
	db := mysql.MustNewDB(provider.LoadMySQLConfig(), gormLog)

	v := bind.NewValidate()
	t := bind.MustNewTrans(v)
	binder := bind.New(v, t, provider.ProvideBindErrCode())

	rdb := redis.New(provider.LoadRedisConfig())
	engine := admin.NewGin(log, db, binder, rdb)

	srv := server.New(provider.LoadAdminServerConfig(), engine, log, provider.ProvideServerOptions()...)
	srv.Run()
}

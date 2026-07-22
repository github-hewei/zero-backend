package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/platform"
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
	log := provider.NewLogger(conn.DB, "platform.log")

	gormLog := gormutil.NewLogger(log)
	db := mysql.MustNewDB(provider.LoadMySQLConfig(), gormLog)

	validate := bind.NewValidate()
	trans := bind.MustNewTrans(validate)
	binder := bind.New(validate, trans, provider.ProvideBindErrCode())

	rdb := redis.New(provider.LoadRedisConfig())
	engine := platform.NewGin(log, db, binder, rdb)

	// 启动服务
	server.New(provider.LoadPlatformServerConfig(), engine, log, provider.ProvideServerOptions()...).Run()
}

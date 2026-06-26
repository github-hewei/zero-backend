package main

import (
	"zero-backend/internal/admin"
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
	cfg := config.New()

	conn, err := mongodb.NewConn(app.NewMongoDBConfig(cfg))
	if err != nil {
		panic(err)
	}
	l := app.ProvideLogger(cfg.Logger, conn.DB)

	gormLog := gormutil.NewLogger(l)
	db, err := mysql.NewDB(app.NewMySQLConfig(cfg), gormLog)
	if err != nil {
		panic(err)
	}

	v := bind.NewValidate()
	t, err := bind.NewTrans(v)
	if err != nil {
		panic(err)
	}
	binder := bind.New(v, t, app.ProvideBindErrCode())

	rdb := redis.New(app.NewRedisConfig(cfg))
	captchaSvc := app.NewCaptchaService(rdb, app.NewCaptchaConfig(cfg))

	srv := server.New(
		app.NewAdminServerConfig(cfg),
		admin.NewGin(l, db, binder,
			app.NewAdminAuthConfig(cfg),
			app.NewAdminCorsConfig(cfg),
			app.NewSettingService(db),
			captchaSvc,
		),
		l,
		app.ProvideServerOptions()...,
	)
	srv.Run()
}

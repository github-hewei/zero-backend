package providers

import (
	"zero-backend/cmd/admin/router"
	apiRouter "zero-backend/cmd/api/router"
	"zero-backend/cmd/cli/command"
	"zero-backend/cmd/cli/runner"
	"zero-backend/cmd/worker/handler"
	workerServer "zero-backend/cmd/worker/server"
	"zero-backend/config"
	"zero-backend/modules/captcha"
	"zero-backend/modules/rbac"
	"zero-backend/modules/setting"
	"zero-backend/modules/user"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/gormutil"
	zeroLogger "github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/queue"
	zeroRedis "github.com/241x/zero-kit/redis"
	baseconfig "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/errcode"
	webserver "github.com/241x/zero-web/server"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// ---------- bind ----------

func ProvideBindErrCode() apperror.Code { return errcode.InvalidInput }

var BindProviderSet = wire.NewSet(
	bind.New, bind.NewValidate, bind.NewTrans, ProvideBindErrCode,
)

// ---------- logger ----------

func ProvideLogger(cfg baseconfig.LoggerConfig, db *mongo.Database) *zeroLogger.ZeroLogger {
	options := []zeroLogger.Option{}
	for _, writer := range cfg.Writers {
		switch writer {
		case "console":
			options = append(options, zeroLogger.WithConsole())
		case "file":
			options = append(options, zeroLogger.WithFileWithConfig(zeroLogger.FileConfig{
				Path: cfg.File.Path, Filename: cfg.File.Filename,
				MaxSize: cfg.File.MaxSize, MaxAge: cfg.File.MaxAge,
				MaxBackups: cfg.File.MaxBackups, Compress: cfg.File.Compress,
				LocalTime: cfg.File.LocalTime,
			}))
		case "mongodb":
			options = append(options, zeroLogger.WithMongo(db))
		}
	}
	level := zeroLogger.Disabled
	switch cfg.Level {
	case "info":
		level = zeroLogger.InfoLevel
	case "debug":
		level = zeroLogger.DebugLevel
	case "warn":
		level = zeroLogger.WarnLevel
	case "error":
		level = zeroLogger.ErrorLevel
	}
	options = append(options, zeroLogger.WithLevel(level))
	return zeroLogger.New(options...)
}

var LoggerProviderSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "Logger"),
	ProvideLogger,
	wire.Bind(new(zeroLogger.Logger), new(*zeroLogger.ZeroLogger)),
)

// ---------- mongodb ----------

func NewMongoDBConfig(cfg *config.Config) mongodb.Config {
	return mongodb.Config{URI: cfg.MongoDB.URI, Database: cfg.MongoDB.Database, Enabled: cfg.MongoDB.Enabled}
}

var MongoDBProviderSet = wire.NewSet(
	NewMongoDBConfig, mongodb.NewConn,
	wire.FieldsOf(new(*mongodb.Conn), "Client", "DB"),
)

// ---------- mysql ----------

func NewMySQLConfig(cfg *config.Config) mysql.Config {
	return mysql.Config{Dsn: cfg.MySQL.Dsn, Prefix: cfg.MySQL.Prefix}
}

var MySQLProviderSet = wire.NewSet(
	NewMySQLConfig, mysql.NewDB, gormutil.NewLogger,
	wire.Bind(new(gormLogger.Interface), new(*gormutil.Logger)),
)

// ---------- redis ----------

func NewRedisConfig(cfg *config.Config) zeroRedis.Config {
	return zeroRedis.Config{Host: cfg.Redis.Host, Port: cfg.Redis.Port, Password: cfg.Redis.Password, DB: cfg.Redis.DB}
}

var RedisProviderSet = wire.NewSet(NewRedisConfig, zeroRedis.New)

// ---------- repository ----------

var RepositoryProviderSet = wire.NewSet(user.NewRepository, rbac.NewRbacApiRepository)

// ---------- server ----------

func ProvideServerOptions() []webserver.Option { return nil }

func NewAdminServerConfig(cfg *config.Config) baseconfig.ServerConfig { return cfg.Admin.Server }
func NewApiServerConfig(cfg *config.Config) baseconfig.ServerConfig   { return cfg.Api.Server }
func NewAdminCorsConfig(cfg *config.Config) baseconfig.CorsConfig     { return cfg.Admin.Cors }
func NewApiCorsConfig(cfg *config.Config) baseconfig.CorsConfig       { return cfg.Api.Cors }

var AdminServerProviderSet = wire.NewSet(
	NewAdminServerConfig, NewAdminCorsConfig, ProvideServerOptions, webserver.New, router.NewGin,
)
var ApiServerProviderSet = wire.NewSet(
	NewApiServerConfig, NewApiCorsConfig, ProvideServerOptions, webserver.New, apiRouter.NewGin,
)

// ---------- service ----------

func NewSettingService(db *gorm.DB) *setting.Service {
	return setting.NewService(setting.NewRepository(db), setting.NewDefaultRepository(db))
}

func NewAdminAuthConfig(cfg *config.Config) config.AdminAuthConfig { return cfg.Admin.Auth }
func NewCaptchaConfig(cfg *config.Config) config.CaptchaConfig     { return cfg.Admin.Captcha }
func NewApiAuthConfig(cfg *config.Config) config.ApiAuthConfig     { return cfg.Api.Auth }

func NewCaptchaService(rdb *goredis.Client, cfg config.CaptchaConfig) *captcha.Service {
	return captcha.NewService(rdb, captcha.Config{Enabled: cfg.Enabled, TTL: cfg.TTL}, "ZAG:CAPTCHA")
}

var AdminServiceProviderSet = wire.NewSet(
	NewAdminAuthConfig, NewCaptchaConfig, NewCaptchaService, NewSettingService,
)
var ApiServiceProviderSet = wire.NewSet(NewApiAuthConfig, NewSettingService)

// ---------- command ----------

var CliCommandProviderSet = wire.NewSet(
	command.NewRootCommand, command.NewMigrateCommand, command.NewQueueCommand,
	command.NewSyncApiCommand, runner.NewSyncApiRunner, queue.NewQueueManager,
)

// ---------- worker ----------

func ProvideRegistry(log zeroLogger.Logger, exampleHandler *handler.ExampleHandler) *handler.Registry {
	registry := handler.NewRegistry(log)
	registry.Register("example", exampleHandler)
	return registry
}

var WorkerProviderSet = wire.NewSet(
	ProvideRegistry, handler.NewExampleHandler, queue.NewQueueManager, workerServer.NewWorkerServer,
)

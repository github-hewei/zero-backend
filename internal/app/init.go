package app

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	baseconfig "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/errcode"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProvideBindErrCode() apperror.Code { return errcode.InvalidInput }

func ProvideLogger(cfg baseconfig.LoggerConfig, db *mongo.Database) *logger.ZeroLogger {
	options := []logger.Option{}
	for _, writer := range cfg.Writers {
		switch writer {
		case "console":
			options = append(options, logger.WithConsole())
		case "file":
			options = append(options, logger.WithFileWithConfig(logger.FileConfig{
				Path: cfg.File.Path, Filename: cfg.File.Filename,
				MaxSize: cfg.File.MaxSize, MaxAge: cfg.File.MaxAge,
				MaxBackups: cfg.File.MaxBackups, Compress: cfg.File.Compress,
				LocalTime: cfg.File.LocalTime,
			}))
		case "mongodb":
			options = append(options, logger.WithMongo(db))
		}
	}
	level := logger.Disabled
	switch cfg.Level {
	case "info":
		level = logger.InfoLevel
	case "debug":
		level = logger.DebugLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	}
	options = append(options, logger.WithLevel(level))
	return logger.New(options...)
}

func NewMongoDBConfig(cfg *config.Config) mongodb.Config {
	return mongodb.Config{URI: cfg.MongoDB.URI, Database: cfg.MongoDB.Database, Enabled: cfg.MongoDB.Enabled}
}

func NewMySQLConfig(cfg *config.Config) mysql.Config {
	return mysql.Config{Dsn: cfg.MySQL.Dsn, Prefix: cfg.MySQL.Prefix}
}

func NewRedisConfig(cfg *config.Config) redis.Config {
	return redis.Config{Host: cfg.Redis.Host, Port: cfg.Redis.Port, Password: cfg.Redis.Password, DB: cfg.Redis.DB}
}

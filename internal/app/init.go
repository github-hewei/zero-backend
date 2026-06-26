package app

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/errcode"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level   string        `mapstructure:"level"`
	Writers []string      `mapstructure:"writers"`
	File    FileLogConfig `mapstructure:"file"`
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string `mapstructure:"path"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"local_time"`
}

func ProvideBindErrCode() apperror.Code { return errcode.InvalidInput }

func LoadLogger(db *mongo.Database) *logger.ZeroLogger {
	var cfg LoggerConfig
	config.UnmarshalKey("logger", &cfg)

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

func LoadMongoConfig() mongodb.Config {
	var c mongodb.Config
	config.UnmarshalKey("mongodb", &c)
	return c
}

func LoadMySQLConfig() mysql.Config {
	var c mysql.Config
	config.UnmarshalKey("mysql", &c)
	return c
}

func LoadRedisConfig() redis.Config {
	var c redis.Config
	config.UnmarshalKey("redis", &c)
	return c
}

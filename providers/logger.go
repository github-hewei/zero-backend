package providers

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/logger"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProvideLogger 提供日志实例
func ProvideLogger(cfg config.LoggerConfig, db *mongo.Database) *logger.ZeroLogger {
	options := []logger.Option{}
	for _, writer := range cfg.Writers {
		switch writer {
		case "console":
			options = append(options, logger.WithConsole())
		case "file":
			options = append(options, logger.WithFileWithConfig(logger.FileConfig{
				Path:       cfg.File.Path,
				Filename:   cfg.File.Filename,
				MaxSize:    cfg.File.MaxSize,
				MaxAge:     cfg.File.MaxAge,
				MaxBackups: cfg.File.MaxBackups,
				Compress:   cfg.File.Compress,
				LocalTime:  cfg.File.LocalTime,
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

// LoggerProviderSet 提供日志依赖集合
var LoggerProviderSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "Logger"),
	ProvideLogger,
	wire.Bind(new(logger.Logger), new(*logger.ZeroLogger)),
)

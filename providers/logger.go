package providers

import (
	"zero-backend/internal/config"
	"zero-backend/internal/logger"
	logger2 "zero-backend/pkg/logger"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProvideLogger 提供日志实例
func ProvideLogger(config *config.Config, db *mongo.Database) *logger.Logger {
	if config == nil {
		return logger.Nop()
	}

	options := []logger.Option{}
	for _, writer := range config.Logger.Writers {
		switch writer {
		case "console":
			options = append(options, logger.WithConsole())
		case "file":
			options = append(options, logger.WithFileWithConfig(logger.FileConfig{
				Path:       config.Logger.File.Path,
				Filename:   config.Logger.File.Filename,
				MaxSize:    config.Logger.File.MaxSize,
				MaxAge:     config.Logger.File.MaxAge,
				MaxBackups: config.Logger.File.MaxBackups,
				Compress:   config.Logger.File.Compress,
				LocalTime:  config.Logger.File.LocalTime,
			}))
		case "mongodb":
			options = append(options, logger.WithMongo(db))
		}
	}

	level := logger2.Disabled
	switch config.Logger.Level {
	case "info":
		level = logger2.InfoLevel
	case "debug":
		level = logger2.DebugLevel
	case "warn":
		level = logger2.WarnLevel
	case "error":
		level = logger2.ErrorLevel
	}

	options = append(options, logger.WithLevel(level))
	return logger.New(options...)
}

// LoggerProviderSet 提供日志依赖集合
var LoggerProviderSet = wire.NewSet(ProvideLogger, wire.Bind(new(logger2.Logger), new(*logger.Logger)))

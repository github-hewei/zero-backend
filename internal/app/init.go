package app

import (
	"strings"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/captcha"
	"zero-backend/internal/modules/setting"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/server"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level   string
	Writers []string
	File    FileLogConfig
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	LocalTime  bool
}

// LoadLogger 加载日志组件。filename 为空时使用配置文件中的默认文件名。
func LoadLogger(db *mongo.Database, filename string) *logger.ZeroLogger {
	var cfg LoggerConfig
	config.UnmarshalKey("logger", &cfg)
	if lv := config.GetString("logger.level"); lv != "" {
		cfg.Level = lv
	}
	if s := config.GetString("logger.writers"); s != "" {
		cfg.Writers = splitComma(s)
	}
	if filename != "" {
		cfg.File.Filename = filename
	}

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

// LoadMongoConfig 加载 MongoDB 配置
func LoadMongoConfig() mongodb.Config {
	return mongodb.Config{
		URI:      config.GetString("mongodb.uri"),
		Database: config.GetString("mongodb.database"),
		Enabled:  config.GetBool("mongodb.enabled"),
	}
}

// LoadMySQLConfig 加载 MySQL 配置
func LoadMySQLConfig() mysql.Config {
	return mysql.Config{
		Dsn:    config.GetString("mysql.dsn"),
		Prefix: config.GetString("mysql.prefix"),
	}
}

// LoadRedisConfig 加载 Redis 配置
func LoadRedisConfig() redis.Config {
	return redis.Config{
		Host:     config.GetString("redis.host"),
		Port:     config.GetInt("redis.port"),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.db"),
	}
}

// splitComma 分割逗号分隔的字符串
func splitComma(s string) []string {
	var res []string
	for item := range strings.SplitSeq(s, ",") {
		if t := strings.TrimSpace(item); t != "" {
			res = append(res, t)
		}
	}
	return res
}

// ProvideBindErrCode 提供绑定组件错误码
func ProvideBindErrCode() apperror.Code { return errcode.InvalidInput }

// ProvideServerOptions 提供服务器选项
func ProvideServerOptions() []server.Option { return nil }

// LoadAdminServerConfig 加载管理后台服务器配置
func LoadAdminServerConfig() server.Config {
	return server.Config{
		Host: config.GetString("admin.server.host"),
		Port: config.GetInt("admin.server.port"),
	}
}

// LoadApiServerConfig 加载 API 服务器配置
func LoadApiServerConfig() server.Config {
	return server.Config{
		Host: config.GetString("api.server.host"),
		Port: config.GetInt("api.server.port"),
	}
}

// LoadPlatformServerConfig 加载平台端服务器配置
func LoadPlatformServerConfig() server.Config {
	return server.Config{
		Host: config.GetString("platform.server.host"),
		Port: config.GetInt("platform.server.port"),
	}
}

// LoadAdminCorsConfig 加载管理后台 CORS 配置
func LoadAdminCorsConfig() middleware.CorsConfig {
	type adapter struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	}
	var a adapter
	config.UnmarshalKey("admin.cors", &a)
	return middleware.CorsConfig{
		AllowOrigins:     a.AllowOrigins,
		AllowMethods:     a.AllowMethods,
		AllowHeaders:     a.AllowHeaders,
		AllowCredentials: a.AllowCredentials,
	}
}

// LoadApiCorsConfig 加载 API CORS 配置
func LoadApiCorsConfig() middleware.CorsConfig {
	type adapter struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	}
	var a adapter
	config.UnmarshalKey("api.cors", &a)
	return middleware.CorsConfig{
		AllowOrigins:     a.AllowOrigins,
		AllowMethods:     a.AllowMethods,
		AllowHeaders:     a.AllowHeaders,
		AllowCredentials: a.AllowCredentials,
	}
}

// LoadCaptchaConfig 加载验证码配置
func LoadCaptchaConfig() captcha.Config {
	type adapter struct {
		Enabled bool `mapstructure:"enabled"`
		TTL     int  `mapstructure:"ttl"`
	}
	var a adapter
	config.UnmarshalKey("admin.captcha", &a)
	return captcha.Config{
		Enabled: a.Enabled,
		TTL:     a.TTL,
	}
}

// NewSettingService 创建设置服务
func NewSettingService(db *gorm.DB) *setting.Service {
	return setting.NewService(setting.NewRepository(db), setting.NewDefaultRepository(db))
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(rdb *goredis.Client, cfg captcha.Config) (*captcha.Service, error) {
	return captcha.NewService(rdb, cfg, "ZAG:CAPTCHA")
}

// Must 泛型辅助：err != nil 时 panic，用于启动阶段集中处理错误。
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

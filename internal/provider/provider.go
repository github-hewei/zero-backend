package provider

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

// LoadLogger 加载日志组件。filename 为空时使用配置文件中的默认文件名。
func LoadLogger(db *mongo.Database, filename string) *logger.ZeroLogger {
	type fileLogConfigAdapter struct {
		Path       string `mapstructure:"path"`
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
		Compress   bool   `mapstructure:"compress"`
		LocalTime  bool   `mapstructure:"local_time"`
	}
	type adapter struct {
		Level   string               `mapstructure:"level"`
		Writers []string             `mapstructure:"writers"`
		File    fileLogConfigAdapter `mapstructure:"file"`
	}
	var a adapter
	config.UnmarshalKey("logger", &a)
	if filename != "" {
		a.File.Filename = filename
	}

	options := []logger.Option{}
	for _, writer := range a.Writers {
		switch writer {
		case "console":
			options = append(options, logger.WithConsole())
		case "file":
			options = append(options, logger.WithFileWithConfig(logger.FileConfig{
				Path:       a.File.Path,
				Filename:   a.File.Filename,
				MaxSize:    a.File.MaxSize,
				MaxAge:     a.File.MaxAge,
				MaxBackups: a.File.MaxBackups,
				Compress:   a.File.Compress,
				LocalTime:  a.File.LocalTime,
			}))
		case "mongodb":
			options = append(options, logger.WithMongo(db))
		}
	}
	level := logger.Disabled
	switch a.Level {
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
	type adapter struct {
		URI      string `mapstructure:"uri"`
		Database string `mapstructure:"database"`
		Enabled  bool   `mapstructure:"enabled"`
	}
	var a adapter
	config.UnmarshalKey("mongodb", &a)
	return mongodb.Config{
		URI:      a.URI,
		Database: a.Database,
		Enabled:  a.Enabled,
	}
}

// LoadMySQLConfig 加载 MySQL 配置
func LoadMySQLConfig() mysql.Config {
	type adapter struct {
		Dsn    string `mapstructure:"dsn"`
		Prefix string `mapstructure:"prefix"`
	}
	var a adapter
	config.UnmarshalKey("mysql", &a)
	return mysql.Config{
		Dsn:    a.Dsn,
		Prefix: a.Prefix,
	}
}

// LoadRedisConfig 加载 Redis 配置
func LoadRedisConfig() redis.Config {
	type adapter struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}
	var a adapter
	config.UnmarshalKey("redis", &a)
	return redis.Config{
		Host:     a.Host,
		Port:     a.Port,
		Password: a.Password,
		DB:       a.DB,
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
func ProvideBindErrCode() apperror.Code {
	return errcode.InvalidInput
}

// ProvideServerOptions 提供服务器选项
func ProvideServerOptions() []server.Option {
	return []server.Option{}
}

// LoadAdminServerConfig 加载管理后台服务器配置
func LoadAdminServerConfig() server.Config {
	type adapter struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	var a adapter
	config.UnmarshalKey("admin.server", &a)
	return server.Config{
		Host: a.Host,
		Port: a.Port,
	}
}

// LoadApiServerConfig 加载 API 服务器配置
func LoadApiServerConfig() server.Config {
	type adapter struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	var a adapter
	config.UnmarshalKey("api.server", &a)
	return server.Config{
		Host: a.Host,
		Port: a.Port,
	}
}

// LoadPlatformServerConfig 加载平台端服务器配置
func LoadPlatformServerConfig() server.Config {
	type adapter struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	var a adapter
	config.UnmarshalKey("platform.server", &a)
	return server.Config{
		Host: a.Host,
		Port: a.Port,
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

// LoadPlatformCorsConfig 加载平台端 CORS 配置
func LoadPlatformCorsConfig() middleware.CorsConfig {
	type adapter struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	}
	var a adapter
	config.UnmarshalKey("platform.cors", &a)
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

// LoadPlatformCaptchaConfig 加载平台端验证码配置
func LoadPlatformCaptchaConfig() captcha.Config {
	type adapter struct {
		Enabled bool `mapstructure:"enabled"`
		TTL     int  `mapstructure:"ttl"`
	}
	var a adapter
	config.UnmarshalKey("platform.captcha", &a)
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
// func Must[T any](v T, err error) T {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return v
// }

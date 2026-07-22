package provider

import (
	"zero-backend/internal/config"
	"zero-backend/internal/modules/captcha"

	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/server"
)

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

// LoadAdminCaptchaConfig 加载管理后台验证码配置
func LoadAdminCaptchaConfig() captcha.Config {
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

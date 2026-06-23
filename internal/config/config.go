package config

import (
	"errors"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	webconfig "github.com/241x/zero-web/config"
)

// 基础配置类型，引用自 zero-web
type (
	ServerConfig   = webconfig.ServerConfig
	CorsConfig     = webconfig.CorsConfig
	AuthConfig     = webconfig.AuthConfig
	LoggerConfig   = webconfig.LoggerConfig
	FileLogConfig  = webconfig.FileLogConfig
	MySQLConfig    = webconfig.MySQLConfig
	RedisConfig    = webconfig.RedisConfig
	MongoDBConfig  = webconfig.MongoDBConfig
)

// Config 全局配置
type Config struct {
	Admin   AdminConfig  `mapstructure:"admin"`
	Api     ApiConfig    `mapstructure:"api"`
	Logger  LoggerConfig `mapstructure:"logger"`
	MySQL   MySQLConfig  `mapstructure:"mysql"`
	Redis   RedisConfig  `mapstructure:"redis"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Enabled bool `mapstructure:"enabled"`
	TTL     int  `mapstructure:"ttl"`
}

// AdminConfig 管理后台配置
type AdminConfig struct {
	Server  ServerConfig    `mapstructure:"server"`
	Auth    AdminAuthConfig `mapstructure:"auth"`
	Cors    CorsConfig      `mapstructure:"cors"`
	Captcha CaptchaConfig   `mapstructure:"captcha"`
}

// AdminAuthConfig 管理端认证配置
type AdminAuthConfig struct {
	AuthConfig  `mapstructure:",squash"`
	SuperUserId int `mapstructure:"super_user_id"`
}

// ApiConfig API模块配置
type ApiConfig struct {
	Server ServerConfig  `mapstructure:"server"`
	Auth   ApiAuthConfig `mapstructure:"auth"`
	Cors   CorsConfig    `mapstructure:"cors"`
}

// ApiAuthConfig API端认证配置
type ApiAuthConfig struct {
	AuthConfig `mapstructure:",squash"`
}

// New 加载配置（支持 YAML + 环境变量覆盖）
func New() *Config {
	return loadConfig()
}

// loadConfig 加载配置
func loadConfig() *Config {
	// 添加配置文件搜索路径（按优先级排序）
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 环境变量覆盖配置
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			log.Println("[WARN] config.yaml not found")
		} else {
			log.Printf("[WARN] read config.yaml failed: %v\n", err)
		}
	}

	// 使用 godotenv 加载 .env 文件并设置环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] .env file not found, skipping\n")
	}

	// 反序列化配置到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("[ERROR] unmarshal config failed: %v\n", err)
		panic(err)
	}

	// 验证配置
	if err := cfg.validate(); err != nil {
		log.Printf("[ERROR] invalid config: %v\n", err)
		panic(err)
	}

	return &cfg
}

// validate 验证配置有效性
func (c *Config) validate() error {
	// 验证 Admin 配置
	if err := c.validateAdmin(); err != nil {
		return err
	}

	// 验证 API 配置
	if err := c.validateApi(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateAdmin() error {
	if c.Admin.Server.Port <= 0 || c.Admin.Server.Port > 65535 {
		return errors.New("invalid admin server port: must be between 1 and 65535")
	}
	if c.Admin.Auth.AccessTokenTtl <= 0 {
		return errors.New("invalid admin access token TTL: must be positive")
	}
	if c.Admin.Auth.RefreshTokenTtl <= 0 {
		return errors.New("invalid admin refresh token TTL: must be positive")
	}
	if c.Admin.Auth.HmacSecret == "" {
		return errors.New("invalid admin HMAC secret: cannot be empty")
	}
	return nil
}

func (c *Config) validateApi() error {
	if c.Api.Server.Port <= 0 || c.Api.Server.Port > 65535 {
		return errors.New("invalid api server port: must be between 1 and 65535")
	}
	if c.Api.Auth.AccessTokenTtl <= 0 {
		return errors.New("invalid api access token TTL: must be positive")
	}
	if c.Api.Auth.RefreshTokenTtl <= 0 {
		return errors.New("invalid api refresh token TTL: must be positive")
	}
	if c.Api.Auth.HmacSecret == "" {
		return errors.New("invalid api HMAC secret: cannot be empty")
	}
	return nil
}

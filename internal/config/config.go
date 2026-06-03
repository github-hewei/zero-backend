package config

import (
	"errors"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	Admin   AdminConfig   `mapstructure:"admin"`
	Api     ApiConfig     `mapstructure:"api"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	Redis   RedisConfig   `mapstructure:"redis"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
}

// AuthConfig 认证配置（Admin 和 Api 共用字段）
type AuthConfig struct {
	RefreshTokenTtl int    `mapstructure:"refresh_token_ttl"`
	AccessTokenTtl  int    `mapstructure:"access_token_ttl"`
	HmacSecret      string `mapstructure:"hmac_secret"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// CorsConfig 跨域配置
type CorsConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// AdminConfig 管理后台配置
type AdminConfig struct {
	Server ServerConfig    `mapstructure:"server"`
	Auth   AdminAuthConfig `mapstructure:"auth"`
	Cors   CorsConfig      `mapstructure:"cors"`
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

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Dsn    string `mapstructure:"dsn"`
	Prefix string `mapstructure:"prefix"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Enabled  bool   `mapstructure:"enabled"`
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

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

// AdminConfig 管理后台配置
type AdminConfig struct {
	Server          ServerConfig `mapstructure:"server"`
	RefreshTokenTtl int          `mapstructure:"refresh_token_ttl"`
	AccessTokenTtl  int          `mapstructure:"access_token_ttl"`
	HmacSecret      string       `mapstructure:"hmac_secret"`
	SuperUserId     int          `mapstructure:"super_user_id"`
	Cors            CorsConfig   `mapstructure:"cors"`
}

// ApiConfig API模块配置
type ApiConfig struct {
	Server          ServerConfig `mapstructure:"server"`
	RefreshTokenTtl int          `mapstructure:"refresh_token_ttl"`
	AccessTokenTtl  int          `mapstructure:"access_token_ttl"`
	HmacSecret      string       `mapstructure:"hmac_secret"`
	Cors            CorsConfig   `mapstructure:"cors"`
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
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
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
	// 验证服务端口
	if c.Admin.Server.Port <= 0 || c.Admin.Server.Port > 65535 {
		return errors.New("invalid server port: must be between 1 and 65535")
	}

	// 验证访问令牌有效期
	if c.Admin.AccessTokenTtl <= 0 {
		return errors.New("invalid access token TTL: must be positive")
	}

	// 验证刷新令牌有效期
	if c.Admin.RefreshTokenTtl <= 0 {
		return errors.New("invalid refresh token TTL: must be positive")
	}

	// 验证JWT密钥
	if c.Admin.HmacSecret == "" {
		return errors.New("invalid HMAC secret: cannot be empty")
	}

	return nil
}

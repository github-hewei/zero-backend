package rbac

import "zero-backend/internal/config"

// Config 模块配置。
type Config struct {
	HmacSecret      string
	AccessTokenTtl  int
	RefreshTokenTtl int
	SuperUserId     int
}

// LoadConfig 从全局配置加载模块配置。
func LoadConfig() Config {
	var c Config
	config.UnmarshalKey("admin.auth", &c)
	return c
}

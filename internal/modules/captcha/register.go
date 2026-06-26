package captcha

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Deps 模块依赖
type Deps struct {
	RDB    *redis.Client
	Binder *bind.Binder
	Config Config
	Prefix string
}

// RegisterWith 注册验证码模块路由（使用外部创建的 Service）
func RegisterWith(rg *gin.RouterGroup, binder *bind.Binder, svc *Service) {
	h := newHandler(binder, svc)
	rg.POST("/captcha/generate", h.Generate)
}

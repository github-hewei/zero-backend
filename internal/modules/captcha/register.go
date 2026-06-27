package captcha

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
)

// Register 注册验证码模块路由（使用外部创建的 Service）
func Register(rg *gin.RouterGroup, binder *bind.Binder, svc *Service) {
	h := newHandler(binder, svc)
	rg.POST("/captcha/generate", h.Generate)
}

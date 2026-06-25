package captcha

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 验证码模块处理器
type Handler struct {
	binder *bind.Binder
	svc    *Service
}

func newHandler(binder *bind.Binder, svc *Service) *Handler {
	return &Handler{binder: binder, svc: svc}
}

func (h *Handler) Generate(c *gin.Context) {
	result, err := h.svc.Generate(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

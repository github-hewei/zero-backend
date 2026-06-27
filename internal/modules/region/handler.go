package region

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 区划模块处理器
type Handler struct {
	binder *bind.Binder
	svc    *Service
}

// newHandler 创建区划模块处理器
func newHandler(binder *bind.Binder, svc *Service) *Handler {
	return &Handler{binder: binder, svc: svc}
}

// Tree 获取区划树
func (h *Handler) Tree(c *gin.Context) {
	result, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

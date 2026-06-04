package controller

import (
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// CaptchaController 验证码控制器
type CaptchaController struct {
	req  *request.Request
	serv *service.CaptchaService
}

// NewCaptchaController 创建验证码控制器实例
func NewCaptchaController(req *request.Request, serv *service.CaptchaService) *CaptchaController {
	return &CaptchaController{
		req:  req,
		serv: serv,
	}
}

// Generate 生成验证码
func (c *CaptchaController) Generate(ctx *gin.Context) {
	result, err := c.serv.Generate(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

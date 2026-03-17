package controller

import (
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingController 设置控制器
type SettingController struct {
	serv *service.SettingService
}

// NewSettingController 创建设置控制器
func NewSettingController(serv *service.SettingService) *SettingController {
	return &SettingController{serv: serv}
}

// QiniuToken 获取七牛上传token
func (c *SettingController) QiniuToken(ctx *gin.Context) {
	result, err := c.serv.QiniuToken(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

package controller

import (
	"zero-backend/internal/response"

	"github.com/gin-gonic/gin"
)

// HealthController 健康控制器
type HealthController struct{}

// NewHealthController 创建健康控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health 检测服务健康接口
func (c *HealthController) Health(ctx *gin.Context) {
	response.Success(ctx, "ok", gin.H{"status": "UP"})
}

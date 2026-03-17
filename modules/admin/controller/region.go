package controller

import (
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// RegionController 区域控制器
type RegionController struct {
	serv *service.RegionService
}

// NewRegionController 创建区域控制器
func NewRegionController(serv *service.RegionService) *RegionController {
	return &RegionController{serv: serv}
}

// Tree 获取省市区数据
func (c *RegionController) Tree(ctx *gin.Context) {
	result, err := c.serv.Regions(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

package controller

import (
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingController 设置控制器
type SettingController struct {
	req  *request.Request
	serv *service.SettingService
}

// NewSettingController 创建设置控制器
func NewSettingController(req *request.Request, serv *service.SettingService) *SettingController {
	return &SettingController{req: req, serv: serv}
}

// List 获取设置列表
func (c *SettingController) List(ctx *gin.Context) {
	req := &dto.SettingListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建设置
func (c *SettingController) Create(ctx *gin.Context) {
	req := &dto.SettingCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新设置
func (c *SettingController) Update(ctx *gin.Context) {
	req := &dto.SettingUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除设置
func (c *SettingController) Delete(ctx *gin.Context) {
	req := &dto.SettingDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

// FormConfigs 获取设置表单配置
func (c *SettingController) FormConfigs(ctx *gin.Context) {
	req := &dto.SettingFormConfigsRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) {
		req.OnlyPlatform = true
	}

	configs, err := c.serv.FormConfigs(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", configs)
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

// SettingDefaultController 默认设置控制器
type SettingDefaultController struct {
	req  *request.Request
	serv *service.SettingDefaultService
}

// NewSettingDefaultController 创建默认设置控制器
func NewSettingDefaultController(req *request.Request, serv *service.SettingDefaultService) *SettingDefaultController {
	return &SettingDefaultController{req: req, serv: serv}
}

// List 获取默认设置列表
func (c *SettingDefaultController) List(ctx *gin.Context) {
	req := &dto.SettingDefaultListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建默认设置
func (c *SettingDefaultController) Create(ctx *gin.Context) {
	req := &dto.SettingDefaultCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新默认设置
func (c *SettingDefaultController) Update(ctx *gin.Context) {
	req := &dto.SettingDefaultUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除默认设置
func (c *SettingDefaultController) Delete(ctx *gin.Context) {
	req := &dto.SettingDefaultDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

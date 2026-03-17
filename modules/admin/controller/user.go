package controller

import (
	"zero-backend/internal/constants"
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	req  *request.Request
	serv *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(req *request.Request, serv *service.UserService) *UserController {
	return &UserController{req: req, serv: serv}
}

// List 获取用户列表
func (c *UserController) List(ctx *gin.Context) {
	req := &dto.UserListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建用户
func (c *UserController) Create(ctx *gin.Context) {
	req := &dto.UserCreateRequest{}
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

// Update 修改用户信息
func (c *UserController) Update(ctx *gin.Context) {
	req := &dto.UserUpdateRequest{}
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

// Delete 删除用户
func (c *UserController) Delete(ctx *gin.Context) {
	req := &dto.UserDeleteRequest{}
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

// GetPointsLogs 获取用户积分记录
func (c *UserController) GetPointsLogs(ctx *gin.Context) {
	req := &dto.UserPointsLogListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.GetPointsLogs(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// ChangePoints 变更用户积分
func (c *UserController) ChangePoints(ctx *gin.Context) {
	req := &dto.UserPointsChangeRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.SourceType = int8(constants.PointsSourceAdmin)
	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.ChangeUserPoints(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "积分变更成功", nil)
}

// Detail 获取用户详情
func (c *UserController) Detail(ctx *gin.Context) {
	req := &dto.UserDetailRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	user, err := c.serv.Detail(ctx.Request.Context(), req.Id)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", user)
}

package controller

import (
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// RbacUserController 用户控制器
type RbacUserController struct {
	req  *request.Request
	serv *service.RbacUserService
}

// NewRbacUserController 创建用户控制器
func NewRbacUserController(req *request.Request, serv *service.RbacUserService) *RbacUserController {
	return &RbacUserController{req: req, serv: serv}
}

// List 用户列表
func (c *RbacUserController) List(ctx *gin.Context) {
	req := &dto.RbacUserListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	result, err := c.serv.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建用户
func (c *RbacUserController) Create(ctx *gin.Context) {
	req := &dto.RbacUserCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	// 如果不是超级管理员，或者未指定企业，则使用当前用户所属企业
	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新用户
func (c *RbacUserController) Update(ctx *gin.Context) {
	req := &dto.RbacUserUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	// 如果不是超级管理员，或者未指定企业，则使用当前用户所属企业
	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除用户
func (c *RbacUserController) Delete(ctx *gin.Context) {
	req := &dto.RbacUserDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	// 如果不是超级管理员，或者未指定企业，则使用当前用户所属企业
	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

// SetRoles 设置用户角色
func (c *RbacUserController) SetRoles(ctx *gin.Context) {
	req := &dto.RbacUserRoleSetRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.SetRoles(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "设置成功", nil)
}

// ResetPassword 重置用户密码
func (c *RbacUserController) ResetPassword(ctx *gin.Context) {
	req := &dto.RbacUserResetPasswordRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	newPassword, err := c.serv.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "重置成功", gin.H{
		"new_password": newPassword,
	})
}

// RbacMenuController 菜单控制器
type RbacMenuController struct {
	req  *request.Request
	serv *service.RbacMenuService
}

// NewRbacMenuController 创建菜单控制器
func NewRbacMenuController(req *request.Request, serv *service.RbacMenuService) *RbacMenuController {
	return &RbacMenuController{req: req, serv: serv}
}

// List 菜单列表
func (c *RbacMenuController) List(ctx *gin.Context) {
	result, err := c.serv.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建菜单
func (c *RbacMenuController) Create(ctx *gin.Context) {
	req := &dto.RbacMenuCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "保存成功", nil)
}

// Update 修改菜单
func (c *RbacMenuController) Update(ctx *gin.Context) {
	req := &dto.RbacMenuUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "保存成功", nil)
}

// Delete 删除菜单
func (c *RbacMenuController) Delete(ctx *gin.Context) {
	req := &dto.RbacMenuDeleteRequest{}
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

// Sync 同步菜单
func (c *RbacMenuController) Sync(ctx *gin.Context) {
	var req []dto.RbacMenuSyncRequest
	if err := c.req.ShouldBindJSONArray(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Sync(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "同步成功", nil)
}

// ApiList 获取菜单关联接口
func (c *RbacMenuController) ApiList(ctx *gin.Context) {
	req := &dto.RbacMenuApiListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.FindApiList(ctx.Request.Context(), req)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "获取成功", result)
}

// ApiSave 保存菜单关联接口
func (c *RbacMenuController) ApiSave(ctx *gin.Context) {
	req := &dto.RbacMenuApiSaveRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.SaveApiList(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "保存成功", nil)
}

// RbacRoleController 角色控制器
type RbacRoleController struct {
	req  *request.Request
	serv *service.RbacRoleService
}

// NewRbacRoleController 创建角色控制器
func NewRbacRoleController(serv *service.RbacRoleService, req *request.Request) *RbacRoleController {
	return &RbacRoleController{req: req, serv: serv}
}

// List 角色列表
func (c *RbacRoleController) List(ctx *gin.Context) {
	req := &dto.RbacRoleListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	result, err := c.serv.FindTreeList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建角色
func (c *RbacRoleController) Create(ctx *gin.Context) {
	req := &dto.RbacRoleCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新角色
func (c *RbacRoleController) Update(ctx *gin.Context) {
	req := &dto.RbacRoleUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除角色
func (c *RbacRoleController) Delete(ctx *gin.Context) {
	req := &dto.RbacRoleDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if !request.IsSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = request.GetStoreId(ctx)
	}

	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

// RbacApiController API权限控制器
type RbacApiController struct {
	req  *request.Request
	serv *service.RbacApiService
}

// NewRbacApiController 创建API权限控制器
func NewRbacApiController(req *request.Request, serv *service.RbacApiService) *RbacApiController {
	return &RbacApiController{req: req, serv: serv}
}

// List 获取API列表
func (c *RbacApiController) List(ctx *gin.Context) {
	result, err := c.serv.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建API
func (c *RbacApiController) Create(ctx *gin.Context) {
	req := &dto.RbacApiCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "保存成功", nil)
}

// Update 更新API
func (c *RbacApiController) Update(ctx *gin.Context) {
	req := &dto.RbacApiUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "保存成功", nil)
}

// Delete 删除API
func (c *RbacApiController) Delete(ctx *gin.Context) {
	req := &dto.RbacApiDeleteRequest{}
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

// RbacStoreController 企业控制器
type RbacStoreController struct {
	req  *request.Request
	serv *service.RbacStoreService
}

// NewRbacStoreController 创建企业控制器
func NewRbacStoreController(req *request.Request, serv *service.RbacStoreService) *RbacStoreController {
	return &RbacStoreController{req: req, serv: serv}
}

// List 获取企业列表
func (c *RbacStoreController) List(ctx *gin.Context) {
	req := &dto.RbacStoreListRequest{}
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

// Create 创建企业
func (c *RbacStoreController) Create(ctx *gin.Context) {
	req := &dto.RbacStoreCreateRequest{}
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

// Delete 删除企业
func (c *RbacStoreController) Delete(ctx *gin.Context) {
	req := &dto.RbacStoreDeleteRequest{}
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

// Update 更新企业信息
func (c *RbacStoreController) Update(ctx *gin.Context) {
	req := &dto.RbacStoreUpdateRequest{}
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

// Recycle 将企业移入回收站
func (c *RbacStoreController) Recycle(ctx *gin.Context) {
	req := &dto.RbacStoreDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Recycle(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "操作成功", nil)
}

// Restore 将企业从回收站恢复
func (c *RbacStoreController) Restore(ctx *gin.Context) {
	req := &dto.RbacStoreDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.Restore(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "操作成功", nil)
}

// SetMenus 设置角色菜单
func (c *RbacRoleController) SetMenus(ctx *gin.Context) {
	req := &dto.RbacRoleMenuSetRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.SetMenus(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "设置成功", nil)
}

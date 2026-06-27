package rbac

import (
	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// handler rbac 模块处理器
type handler struct {
	binder    *bind.Binder
	authServ  *AuthService
	authCfg   Config
	menuServ  *RbacMenuService
	apiServ   *RbacApiService
	roleServ  *RbacRoleService
	userServ  *RbacUserService
	storeServ *RbacStoreService
	authMid   *AuthMiddleware
}

// newHandler 创建处理器
func newHandler(
	binder *bind.Binder,
	authServ *AuthService,
	authCfg Config,
	menuServ *RbacMenuService,
	apiServ *RbacApiService,
	roleServ *RbacRoleService,
	userServ *RbacUserService,
	storeServ *RbacStoreService,
	authMid *AuthMiddleware,
) *handler {
	return &handler{
		binder:    binder,
		authServ:  authServ,
		authCfg:   authCfg,
		menuServ:  menuServ,
		apiServ:   apiServ,
		roleServ:  roleServ,
		userServ:  userServ,
		storeServ: storeServ,
		authMid:   authMid,
	}
}

// isSuperUser 是否超管
func isSuperUser(ctx *gin.Context) bool {
	user, ok := ctxkeys.User(ctx.Request.Context()).(*RbacUser)
	return ok && user.SU
}

// storeID 获取当前请求的店铺ID
func storeID(ctx *gin.Context) uint32 {
	return ctxkeys.StoreID(ctx.Request.Context())
}

// login 系统登录
// @Summary 系统登录
// @Tags 认证
// @Param body body AuthLoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=AdminLoginResponse}
// @Router /login [post]
func (h *handler) login(ctx *gin.Context) {
	req := &AuthLoginRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	result, refreshToken, err := h.authServ.Login(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	ctx.SetCookie("token", refreshToken, h.authCfg.RefreshTokenTtl, "/api/refresh-token", "", true, true)
	response.Success(ctx, "请求成功", result)
}

// logout 系统登出
// @Summary 系统登出
// @Tags 认证
// @Success 200 {object} response.Response
// @Router /logout [post]
func (h *handler) logout(ctx *gin.Context) {
	response.Success(ctx, "操作成功", nil)
}

// refreshToken 刷新令牌
// @Summary 刷新令牌
// @Tags 认证
// @Success 200 {object} response.Response{data=AdminLoginResponse}
// @Router /refresh-token [post]
func (h *handler) refreshToken(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("token")
	if err != nil {
		response.Error(ctx, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("读取Cookie失败")))
		return
	}
	result, err := h.authServ.RefreshToken(ctx.Request.Context(), cookie.Value)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// changePassword 修改密码
// @Summary 修改密码
// @Tags 认证
// @Param body body ChangePasswordRequest true "修改密码参数"
// @Success 200 {object} response.Response
// @Router /change-password [post]
func (h *handler) changePassword(ctx *gin.Context) {
	req := &ChangePasswordRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.authServ.ChangePassword(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "密码修改成功", nil)
}

// permissions 获取用户权限菜单
// @Summary 获取用户权限菜单
// @Tags 认证
// @Param body body AuthGetPermissionsRequest true "查询参数"
// @Success 200 {object} response.Response{data=[]RbacMenu}
// @Router /permissions [post]
func (h *handler) permissions(ctx *gin.Context) {
	req := &AuthGetPermissionsRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := h.authServ.GetPermissions(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// menuList 获取菜单列表
// @Summary 获取菜单列表
// @Tags RBAC菜单管理
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/menu/list [post]
func (h *handler) menuList(ctx *gin.Context) {
	result, err := h.menuServ.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// menuCreate 创建菜单
// @Summary 创建菜单
// @Tags RBAC菜单管理
// @Param body body RbacMenuCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /rbac/menu/create [post]
func (h *handler) menuCreate(ctx *gin.Context) {
	req := &RbacMenuCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.menuServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// menuUpdate 更新菜单
// @Summary 更新菜单
// @Tags RBAC菜单管理
// @Param body body RbacMenuUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /rbac/menu/update [post]
func (h *handler) menuUpdate(ctx *gin.Context) {
	req := &RbacMenuUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.menuServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// menuDelete 删除菜单
// @Summary 删除菜单
// @Tags RBAC菜单管理
// @Param body body RbacMenuDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /rbac/menu/delete [post]
func (h *handler) menuDelete(ctx *gin.Context) {
	req := &RbacMenuDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.menuServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// menuSync 同步菜单
// @Summary 同步菜单
// @Tags RBAC菜单管理
// @Param body body []RbacMenuSyncRequest true "同步参数"
// @Success 200 {object} response.Response
// @Router /rbac/menu/sync [post]
func (h *handler) menuSync(ctx *gin.Context) {
	var req []RbacMenuSyncRequest
	if err := h.binder.ShouldBindJSONArray(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.menuServ.Sync(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "同步成功", nil)
}

// menuApiList 获取菜单关联的接口列表
// @Summary 获取菜单关联的接口列表
// @Tags RBAC菜单管理
// @Param body body RbacMenuApiListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/menu/api/list [post]
func (h *handler) menuApiList(ctx *gin.Context) {
	req := &RbacMenuApiListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := h.menuServ.FindApiList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "获取成功", result)
}

// menuApiSave 保存菜单关联的接口列表
// @Summary 保存菜单关联的接口列表
// @Tags RBAC菜单管理
// @Param body body RbacMenuApiSaveRequest true "保存参数"
// @Success 200 {object} response.Response
// @Router /rbac/menu/api/save [post]
func (h *handler) menuApiSave(ctx *gin.Context) {
	req := &RbacMenuApiSaveRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.menuServ.SaveApiList(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// apiList 获取接口列表
// @Summary 获取接口列表
// @Tags RBAC接口管理
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/api/list [post]
func (h *handler) apiList(ctx *gin.Context) {
	result, err := h.apiServ.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// apiCreate 创建接口
// @Summary 创建接口
// @Tags RBAC接口管理
// @Param body body RbacApiCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /rbac/api/create [post]
func (h *handler) apiCreate(ctx *gin.Context) {
	req := &RbacApiCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.apiServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// apiUpdate 更新接口
// @Summary 更新接口
// @Tags RBAC接口管理
// @Param body body RbacApiUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /rbac/api/update [post]
func (h *handler) apiUpdate(ctx *gin.Context) {
	req := &RbacApiUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.apiServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// apiDelete 删除接口
// @Summary 删除接口
// @Tags RBAC接口管理
// @Param body body RbacApiDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /rbac/api/delete [post]
func (h *handler) apiDelete(ctx *gin.Context) {
	req := &RbacApiDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.apiServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// storeList 获取企业列表
// @Summary 获取企业列表
// @Tags RBAC企业管理
// @Param body body RbacStoreListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/store/list [post]
func (h *handler) storeList(ctx *gin.Context) {
	req := &RbacStoreListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := h.storeServ.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// storeCreate 创建企业
// @Summary 创建企业
// @Tags RBAC企业管理
// @Param body body RbacStoreCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /rbac/store/create [post]
func (h *handler) storeCreate(ctx *gin.Context) {
	req := &RbacStoreCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.storeServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建成功", nil)
}

// storeUpdate 更新企业
// @Summary 更新企业
// @Tags RBAC企业管理
// @Param body body RbacStoreUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /rbac/store/update [post]
func (h *handler) storeUpdate(ctx *gin.Context) {
	req := &RbacStoreUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.storeServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新成功", nil)
}

// storeDelete 删除企业
// @Summary 删除企业
// @Tags RBAC企业管理
// @Param body body RbacStoreDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /rbac/store/delete [post]
func (h *handler) storeDelete(ctx *gin.Context) {
	req := &RbacStoreDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.storeServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// storeRecycle 回收企业
// @Summary 回收企业
// @Tags RBAC企业管理
// @Param body body RbacStoreDeleteRequest true "回收参数"
// @Success 200 {object} response.Response
// @Router /rbac/store/recycle [post]
func (h *handler) storeRecycle(ctx *gin.Context) {
	req := &RbacStoreDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.storeServ.Recycle(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "操作成功", nil)
}

// storeRestore 恢复企业
// @Summary 恢复企业
// @Tags RBAC企业管理
// @Param body body RbacStoreDeleteRequest true "恢复参数"
// @Success 200 {object} response.Response
// @Router /rbac/store/restore [post]
func (h *handler) storeRestore(ctx *gin.Context) {
	req := &RbacStoreDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.storeServ.Restore(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "操作成功", nil)
}

// roleList 获取角色列表
// @Summary 获取角色列表
// @Tags RBAC角色管理
// @Param body body RbacRoleListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/role/list [post]
func (h *handler) roleList(ctx *gin.Context) {
	req := &RbacRoleListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	result, err := h.roleServ.FindTreeList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// roleCreate 创建角色
// @Summary 创建角色
// @Tags RBAC角色管理
// @Param body body RbacRoleCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /rbac/role/create [post]
func (h *handler) roleCreate(ctx *gin.Context) {
	req := &RbacRoleCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.roleServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建成功", nil)
}

// roleUpdate 更新角色
// @Summary 更新角色
// @Tags RBAC角色管理
// @Param body body RbacRoleUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /rbac/role/update [post]
func (h *handler) roleUpdate(ctx *gin.Context) {
	req := &RbacRoleUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.roleServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新成功", nil)
}

// roleDelete 删除角色
// @Summary 删除角色
// @Tags RBAC角色管理
// @Param body body RbacRoleDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /rbac/role/delete [post]
func (h *handler) roleDelete(ctx *gin.Context) {
	req := &RbacRoleDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.roleServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// roleSetMenus 设置角色菜单
// @Summary 设置角色菜单
// @Tags RBAC角色管理
// @Param body body RbacRoleMenuSetRequest true "设置参数"
// @Success 200 {object} response.Response
// @Router /rbac/role/set-menus [post]
func (h *handler) roleSetMenus(ctx *gin.Context) {
	req := &RbacRoleMenuSetRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.roleServ.SetMenus(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "设置成功", nil)
}

// userList 获取后台用户列表
// @Summary 获取后台用户列表
// @Tags RBAC用户管理
// @Param body body RbacUserListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /rbac/user/list [post]
func (h *handler) userList(ctx *gin.Context) {
	req := &RbacUserListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	result, err := h.userServ.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// userCreate 创建后台用户
// @Summary 创建后台用户
// @Tags RBAC用户管理
// @Param body body RbacUserCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /rbac/user/create [post]
func (h *handler) userCreate(ctx *gin.Context) {
	req := &RbacUserCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.userServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建成功", nil)
}

// userUpdate 更新后台用户
// @Summary 更新后台用户
// @Tags RBAC用户管理
// @Param body body RbacUserUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /rbac/user/update [post]
func (h *handler) userUpdate(ctx *gin.Context) {
	req := &RbacUserUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.userServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新成功", nil)
}

// userDelete 删除后台用户
// @Summary 删除后台用户
// @Tags RBAC用户管理
// @Param body body RbacUserDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /rbac/user/delete [post]
func (h *handler) userDelete(ctx *gin.Context) {
	req := &RbacUserDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.userServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// userSetRoles 设置用户角色
// @Summary 设置用户角色
// @Tags RBAC用户管理
// @Param body body RbacUserRoleSetRequest true "设置参数"
// @Success 200 {object} response.Response
// @Router /rbac/user/set-roles [post]
func (h *handler) userSetRoles(ctx *gin.Context) {
	req := &RbacUserRoleSetRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	if err := h.userServ.SetRoles(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "设置成功", nil)
}

// userResetPassword 重置用户密码
// @Summary 重置用户密码
// @Tags RBAC用户管理
// @Param body body RbacUserResetPasswordRequest true "重置参数"
// @Success 200 {object} response.Response
// @Router /rbac/user/reset-password [post]
func (h *handler) userResetPassword(ctx *gin.Context) {
	req := &RbacUserResetPasswordRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if !isSuperUser(ctx) || req.StoreId == 0 {
		req.StoreId = storeID(ctx)
	}
	newPassword, err := h.userServ.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "重置成功", gin.H{
		"new_password": newPassword,
	})
}

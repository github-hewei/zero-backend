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

func isSuperUser(ctx *gin.Context) bool {
	user, ok := ctxkeys.User(ctx.Request.Context()).(*RbacUser)
	return ok && user.SU
}

func storeID(ctx *gin.Context) uint32 {
	return ctxkeys.StoreID(ctx.Request.Context())
}

// --- auth ---

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

// @Summary 系统登出
// @Tags 认证
// @Success 200 {object} response.Response
// @Router /logout [post]
func (h *handler) logout(ctx *gin.Context) {
	response.Success(ctx, "操作成功", nil)
}

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

// --- rbac menu ---

func (h *handler) menuList(ctx *gin.Context) {
	result, err := h.menuServ.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

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

// --- rbac api ---

func (h *handler) apiList(ctx *gin.Context) {
	result, err := h.apiServ.FindTreeList(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

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

// --- rbac store ---

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

// --- rbac role ---

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

// --- rbac user ---

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

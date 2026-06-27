package platform_user

import (
	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// handler 平台模块处理器
type handler struct {
	binder   *bind.Binder
	authServ *AuthService
	authCfg  Config
	userServ *PlatformUserService
}

// newHandler 创建处理器
func newHandler(
	binder *bind.Binder,
	authServ *AuthService,
	authCfg Config,
	userServ *PlatformUserService,
) *handler {
	return &handler{
		binder:   binder,
		authServ: authServ,
		authCfg:  authCfg,
		userServ: userServ,
	}
}

// login 平台登录
func (h *handler) login(ctx *gin.Context) {
	req := &PlatformLoginRequest{}
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

// logout 平台登出
func (h *handler) logout(ctx *gin.Context) {
	response.Success(ctx, "操作成功", nil)
}

// refreshToken 刷新令牌
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

// userList 平台用户列表
func (h *handler) userList(ctx *gin.Context) {
	req := &PlatformUserListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := h.userServ.FindList(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// userCreate 创建平台用户
func (h *handler) userCreate(ctx *gin.Context) {
	req := &PlatformUserCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.userServ.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// userUpdate 更新平台用户
func (h *handler) userUpdate(ctx *gin.Context) {
	req := &PlatformUserUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.userServ.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "保存成功", nil)
}

// userDelete 删除平台用户
func (h *handler) userDelete(ctx *gin.Context) {
	req := &PlatformUserDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.userServ.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// userResetPassword 重置平台用户密码
func (h *handler) userResetPassword(ctx *gin.Context) {
	req := &PlatformUserResetPasswordRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := h.userServ.ResetPassword(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "重置成功", nil)
}

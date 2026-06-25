package controller

import (
	"zero-backend/internal/config"
	"zero-backend/modules/rbac"
	"zero-backend/modules/api/service"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// AuthController 权限控制器
type AuthController struct {
	req  *bind.Binder
	serv *service.AuthService
	cfg  config.ApiAuthConfig
}

// NewAuthController 创建权限控制器实例
func NewAuthController(req *bind.Binder, serv *service.AuthService, cfg config.ApiAuthConfig) *AuthController {
	return &AuthController{
		req:  req,
		serv: serv,
		cfg:  cfg,
	}
}

// Login 系统登录
func (c *AuthController) Login(ctx *gin.Context) {
	req := &rbac.AuthLoginRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, refreshToken, err := c.serv.Login(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	ctx.SetCookie("token", refreshToken, c.cfg.RefreshTokenTtl, "/api/refresh-token", "", true, true)
	response.Success(ctx, "请求成功", result)
}

// Logout 系统登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 实现退出登录功能，前端丢弃jwt即可。(后端也可以把jwt写入黑名单并拒绝后续请求)
	response.Success(ctx, "操作成功", nil)
}

// RefreshToken 刷新AccessToken Token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("token")
	if err != nil {
		response.Error(ctx, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("读取Cookie失败")))
		return
	}

	result, err := c.serv.RefreshToken(ctx.Request.Context(), cookie.Value)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// ChangePassword 修改密码
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	req := &rbac.ChangePasswordRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := c.serv.ChangePassword(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "密码修改成功", nil)
}

package controller

import (
	"zero-backend/internal/apperror"
	"zero-backend/internal/config"
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/modules/admin/service"

	"github.com/gin-gonic/gin"
)

// AuthController 权限控制器
type AuthController struct {
	req  *request.Request
	serv *service.AuthService
	cfg  *config.Config
}

// NewAuthController 创建权限控制器实例
func NewAuthController(req *request.Request, serv *service.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		req:  req,
		serv: serv,
		cfg:  cfg,
	}
}

// Login 系统登录
func (c *AuthController) Login(ctx *gin.Context) {
	req := &dto.AuthLoginRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, refreshToken, err := c.serv.Login(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.SetCookie(ctx, "token", refreshToken, c.cfg.Admin.RefreshTokenTtl, "/api/refresh-token")
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
		response.Error(ctx, apperror.NewSystemError(err, "登录已过期，请重新登录"))
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
	req := &dto.ChangePasswordRequest{}
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

// GetPermissions 获取用户菜单权限
func (c *AuthController) GetPermissions(ctx *gin.Context) {
	req := &dto.AuthGetPermissionsRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	result, err := c.serv.GetPermissions(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

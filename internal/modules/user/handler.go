package user

import (
	"zero-backend/internal/modules/rbac"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 用户模块处理器
type Handler struct {
	binder   *bind.Binder
	svc      *Service
	authServ *AuthService
	authCfg  Config
}

// NewHandler 创建用户模块处理器
func newHandler(binder *bind.Binder, svc *Service) *Handler {
	return &Handler{binder: binder, svc: svc}
}

// NewAuthHandler 创建用户认证模块处理器
func newAuthHandler(binder *bind.Binder, svc *Service, authServ *AuthService, authCfg Config) *Handler {
	return &Handler{binder: binder, svc: svc, authServ: authServ, authCfg: authCfg}
}

// List 获取用户列表
// @Summary 获取用户列表
// @Tags 用户管理
// @Param body body ListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /user/user/list [post]
func (h *Handler) List(c *gin.Context) {
	req := &ListRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	result, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Param body body CreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /user/user/create [post]
func (h *Handler) Create(c *gin.Context) {
	req := &CreateRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	if err := h.svc.Create(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "创建成功", nil)
}

// Update 更新用户
// @Summary 更新用户
// @Tags 用户管理
// @Param body body UpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /user/user/update [post]
func (h *Handler) Update(c *gin.Context) {
	req := &UpdateRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	if err := h.svc.Update(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "更新成功", nil)
}

// Delete 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Param body body DeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /user/user/delete [post]
func (h *Handler) Delete(c *gin.Context) {
	req := &DeleteRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	if err := h.svc.Delete(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "删除成功", nil)
}

// GetPointsLogs 获取用户积分记录
// @Summary 获取用户积分记录
// @Tags 用户管理
// @Param body body PointsLogListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /user/points/logs [post]
func (h *Handler) GetPointsLogs(c *gin.Context) {
	req := &PointsLogListRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	result, err := h.svc.GetPointsLogs(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// ChangePoints 用户积分变更
// @Summary 用户积分变更
// @Tags 用户管理
// @Param body body PointsChangeRequest true "变更参数"
// @Success 200 {object} response.Response
// @Router /user/points/change [post]
func (h *Handler) ChangePoints(c *gin.Context) {
	req := &PointsChangeRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(c.Request.Context())
	if err := h.svc.ChangePoints(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "积分变动成功", nil)
}

// Detail 获取用户详情
// @Summary 获取用户详情
// @Tags 用户管理
// @Param body body DetailRequest true "查询参数"
// @Success 200 {object} response.Response
// @Router /user/user/detail [post]
func (h *Handler) Detail(c *gin.Context) {
	req := &DetailRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.svc.Detail(c.Request.Context(), req.Id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// login 用户登录
// @Summary 用户登录
// @Tags 用户认证
// @Param body body rbac.AuthLoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=UserLoginResponse}
// @Router /login [post]
func (h *Handler) login(c *gin.Context) {
	req := &rbac.AuthLoginRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	result, refreshToken, err := h.authServ.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.SetCookie("token", refreshToken, h.authCfg.RefreshTokenTtl, "/api/refresh-token", "", true, true)
	response.Success(c, "请求成功", result)
}

// logout 用户登出
// @Summary 用户登出
// @Tags 用户认证
// @Success 200 {object} response.Response
// @Router /logout [post]
func (h *Handler) logout(c *gin.Context) {
	response.Success(c, "操作成功", nil)
}

// refreshToken 用户刷新令牌
// @Summary 用户刷新令牌
// @Tags 用户认证
// @Success 200 {object} response.Response{data=UserLoginResponse}
// @Router /refresh-token [post]
func (h *Handler) refreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		response.Error(c, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("读取Cookie失败")))
		return
	}
	result, err := h.authServ.RefreshToken(c.Request.Context(), cookie.Value)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// changePassword 用户修改密码
// @Summary 用户修改密码
// @Tags 用户认证
// @Param body body rbac.ChangePasswordRequest true "修改密码参数"
// @Success 200 {object} response.Response
// @Router /change-password [post]
func (h *Handler) changePassword(c *gin.Context) {
	req := &rbac.ChangePasswordRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.authServ.ChangePassword(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "密码修改成功", nil)
}

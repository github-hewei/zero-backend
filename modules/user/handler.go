package user

import (
	"zero-backend/internal/config"
	"zero-backend/modules/rbac"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 用户模块处理器
type Handler struct {
	binder   *bind.Binder
	svc      *Service
	authServ *AuthService
	authCfg  config.ApiAuthConfig
}

func newHandler(binder *bind.Binder, svc *Service) *Handler {
	return &Handler{binder: binder, svc: svc}
}

func newAuthHandler(binder *bind.Binder, svc *Service, authServ *AuthService, authCfg config.ApiAuthConfig) *Handler {
	return &Handler{binder: binder, svc: svc, authServ: authServ, authCfg: authCfg}
}

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

func (h *Handler) logout(c *gin.Context) {
	response.Success(c, "操作成功", nil)
}

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

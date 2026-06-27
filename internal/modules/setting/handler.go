package setting

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 设置模块处理器
type Handler struct {
	binder     *bind.Binder
	svc        *Service
	defaultSvc *DefaultService
}

// newHandler 创建设置模块处理器
func newHandler(binder *bind.Binder, svc *Service, defaultSvc *DefaultService) *Handler {
	return &Handler{binder: binder, svc: svc, defaultSvc: defaultSvc}
}

// List 获取设置列表
// @Summary 获取设置列表
// @Tags 设置管理
// @Param body body ListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /setting/list [post]
func (h *Handler) List(c *gin.Context) {
	req := &ListRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.svc.FindList(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// Create 创建设置
// @Summary 创建设置
// @Tags 设置管理
// @Param body body CreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /setting/create [post]
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

// Update 更新设置
// @Summary 更新设置
// @Tags 设置管理
// @Param body body UpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /setting/update [post]
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

// Delete 删除设置
// @Summary 删除设置
// @Tags 设置管理
// @Param body body DeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /setting/delete [post]
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

// FormConfigs 获取表单配置
// @Summary 获取表单配置
// @Tags 设置管理
// @Param body body FormConfigsRequest true "查询参数"
// @Success 200 {object} response.Response
// @Router /setting/form-configs [post]
func (h *Handler) FormConfigs(c *gin.Context) {
	req := &FormConfigsRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.svc.FormConfigs(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// QiniuToken 获取七牛上传Token
// @Summary 获取七牛上传Token
// @Tags 设置管理
// @Success 200 {object} response.Response{data=QiniuTokenResponse}
// @Router /setting/qiniu-token [post]
func (h *Handler) QiniuToken(c *gin.Context) {
	result, err := h.svc.QiniuToken(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// DefaultList 获取默认设置列表
// @Summary 获取默认设置列表
// @Tags 默认设置管理
// @Param body body DefaultListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /setting/default/list [post]
func (h *Handler) DefaultList(c *gin.Context) {
	req := &DefaultListRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.defaultSvc.FindList(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "请求成功", result)
}

// DefaultCreate 创建默认设置
// @Summary 创建默认设置
// @Tags 默认设置管理
// @Param body body DefaultCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /setting/default/create [post]
func (h *Handler) DefaultCreate(c *gin.Context) {
	req := &DefaultCreateRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.defaultSvc.Create(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "创建成功", nil)
}

// DefaultUpdate 更新默认设置
// @Summary 更新默认设置
// @Tags 默认设置管理
// @Param body body DefaultUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /setting/default/update [post]
func (h *Handler) DefaultUpdate(c *gin.Context) {
	req := &DefaultUpdateRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.defaultSvc.Update(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "更新成功", nil)
}

// DefaultDelete 删除默认设置
// @Summary 删除默认设置
// @Tags 默认设置管理
// @Param body body DefaultDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /setting/default/delete [post]
func (h *Handler) DefaultDelete(c *gin.Context) {
	req := &DefaultDeleteRequest{}
	if err := h.binder.ShouldBindJSON(c, req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.defaultSvc.Delete(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "删除成功", nil)
}

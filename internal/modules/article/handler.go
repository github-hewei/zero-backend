package article

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Handler 文章模块处理器
type Handler struct {
	binder      *bind.Binder
	categorySvc *CategoryService
	articleSvc  *Service
}

// NewHandler 创建处理器
func NewHandler(binder *bind.Binder, categorySvc *CategoryService, articleSvc *Service) *Handler {
	return &Handler{binder: binder, categorySvc: categorySvc, articleSvc: articleSvc}
}

// ListCategory 获取文章分类列表
// @Summary 获取文章分类列表
// @Tags 文章分类管理
// @Param body body CategoryListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /article/category/list [post]
func (h *Handler) ListCategory(ctx *gin.Context) {
	req := &CategoryListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	result, err := h.categorySvc.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// CreateCategory 创建文章分类
// @Summary 创建文章分类
// @Tags 文章分类管理
// @Param body body CategoryCreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /article/category/create [post]
func (h *Handler) CreateCategory(ctx *gin.Context) {
	req := &CategoryCreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.categorySvc.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建成功", nil)
}

// UpdateCategory 更新文章分类
// @Summary 更新文章分类
// @Tags 文章分类管理
// @Param body body CategoryUpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /article/category/update [post]
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	req := &CategoryUpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.categorySvc.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新成功", nil)
}

// DeleteCategory 删除文章分类
// @Summary 删除文章分类
// @Tags 文章分类管理
// @Param body body CategoryDeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /article/category/delete [post]
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	req := &CategoryDeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.categorySvc.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

// List 获取文章列表
// @Summary 获取文章列表
// @Tags 文章管理
// @Param body body ListRequest true "查询参数"
// @Success 200 {object} response.Response{data=ListResult}
// @Router /article/article/list [post]
func (h *Handler) List(ctx *gin.Context) {
	req := &ListRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	result, err := h.articleSvc.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "请求成功", result)
}

// Create 创建文章
// @Summary 创建文章
// @Tags 文章管理
// @Param body body CreateRequest true "创建参数"
// @Success 200 {object} response.Response
// @Router /article/article/create [post]
func (h *Handler) Create(ctx *gin.Context) {
	req := &CreateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.articleSvc.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建成功", nil)
}

// Update 更新文章
// @Summary 更新文章
// @Tags 文章管理
// @Param body body UpdateRequest true "更新参数"
// @Success 200 {object} response.Response
// @Router /article/article/update [post]
func (h *Handler) Update(ctx *gin.Context) {
	req := &UpdateRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.articleSvc.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新成功", nil)
}

// Delete 删除文章
// @Summary 删除文章
// @Tags 文章管理
// @Param body body DeleteRequest true "删除参数"
// @Success 200 {object} response.Response
// @Router /article/article/delete [post]
func (h *Handler) Delete(ctx *gin.Context) {
	req := &DeleteRequest{}
	if err := h.binder.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}
	req.StoreId = ctxkeys.StoreID(ctx.Request.Context())
	if err := h.articleSvc.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除成功", nil)
}

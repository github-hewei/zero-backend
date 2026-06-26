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

func NewHandler(binder *bind.Binder, categorySvc *CategoryService, articleSvc *Service) *Handler {
	return &Handler{binder: binder, categorySvc: categorySvc, articleSvc: articleSvc}
}

// --- category ---

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

// --- article ---

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

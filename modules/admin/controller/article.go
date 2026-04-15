package controller

import (
	"zero-backend/internal/dto"
	"zero-backend/internal/request"
	"zero-backend/internal/response"
	"zero-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ArticleCategoryController 文章分类控制器
type ArticleCategoryController struct {
	req  *request.Request
	serv *service.ArticleCategoryService
}

// NewArticleCategoryController 创建文章分类控制器
func NewArticleCategoryController(req *request.Request, serv *service.ArticleCategoryService) *ArticleCategoryController {
	return &ArticleCategoryController{req: req, serv: serv}
}

// List 获取文章分类列表
func (c *ArticleCategoryController) List(ctx *gin.Context) {
	req := &dto.ArticleCategoryListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	result, err := c.serv.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建文章分类
func (c *ArticleCategoryController) Create(ctx *gin.Context) {
	req := &dto.ArticleCategoryCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新文章分类
func (c *ArticleCategoryController) Update(ctx *gin.Context) {
	req := &dto.ArticleCategoryUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除文章分类
func (c *ArticleCategoryController) Delete(ctx *gin.Context) {
	req := &dto.ArticleCategoryDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

// ArticleController 文章控制器
type ArticleController struct {
	req  *request.Request
	serv *service.ArticleService
}

// NewArticleController 创建文章控制器
func NewArticleController(req *request.Request, serv *service.ArticleService) *ArticleController {
	return &ArticleController{req: req, serv: serv}
}

// List 获取文章列表
func (c *ArticleController) List(ctx *gin.Context) {
	req := &dto.ArticleListRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	result, err := c.serv.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", result)
}

// Create 创建文章
func (c *ArticleController) Create(ctx *gin.Context) {
	req := &dto.ArticleCreateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Create(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建成功", nil)
}

// Update 更新文章
func (c *ArticleController) Update(ctx *gin.Context) {
	req := &dto.ArticleUpdateRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Update(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新成功", nil)
}

// Delete 删除文章
func (c *ArticleController) Delete(ctx *gin.Context) {
	req := &dto.ArticleDeleteRequest{}
	if err := c.req.ShouldBindJSON(ctx, req); err != nil {
		response.Error(ctx, err)
		return
	}

	req.StoreId = request.GetStoreId(ctx)
	if err := c.serv.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "删除成功", nil)
}

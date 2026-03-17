package service

import (
	"context"
	"zero-backend/internal/apperror"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
)

// ArticleCategoryService 文章分类业务逻辑层
type ArticleCategoryService struct {
	repo *repository.ArticleCategoryRepository
}

// NewArticleCategoryService 创建文章分类业务逻辑层
func NewArticleCategoryService(repo *repository.ArticleCategoryRepository) *ArticleCategoryService {
	return &ArticleCategoryService{
		repo: repo,
	}
}

// List 获取文章分类列表
func (s *ArticleCategoryService) List(ctx context.Context, req *dto.ArticleCategoryListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.ArticleCategory{},
		Total: 0,
	}

	filter := &repository.ArticleCategoryFilter{
		Name:   req.Name,
		Status: req.Status,
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文章分类总数失败")
	}

	result.Total = total
	if total == 0 {
		return result, nil
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "sort", Sort: "asc"},
		{Field: "id", Sort: "desc"},
	}

	list, err := s.repo.FindAll(ctx, filter, pagination, orders)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文章分类列表失败")
	}

	result.List = list
	return result, nil
}

// Create 创建文章分类
func (s *ArticleCategoryService) Create(ctx context.Context, req *dto.ArticleCategoryCreateRequest) error {
	category := &model.ArticleCategory{
		Name:    req.Name,
		Status:  req.Status,
		Sort:    req.Sort,
		StoreId: req.StoreId,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return apperror.NewSystemError(err, "创建文章分类失败")
	}

	return nil
}

// Update 更新文章分类
func (s *ArticleCategoryService) Update(ctx context.Context, req *dto.ArticleCategoryUpdateRequest) error {
	category, err := s.repo.FindOne(ctx, req.Id)
	if err != nil {
		return apperror.NewSystemError(err, "查询文章分类失败")
	}

	if category.ID == 0 {
		return apperror.NewUserError("文章分类不存在")
	}

	updateData := map[string]any{
		"name":   req.Name,
		"status": req.Status,
		"sort":   req.Sort,
	}

	if err := s.repo.Updates(ctx, category, updateData); err != nil {
		return apperror.NewSystemError(err, "更新文章分类失败")
	}

	return nil
}

// Delete 删除文章分类
func (s *ArticleCategoryService) Delete(ctx context.Context, req *dto.ArticleCategoryDeleteRequest) error {
	category, err := s.repo.FindOne(ctx, req.Id)
	if err != nil {
		return apperror.NewSystemError(err, "查询文章分类失败")
	}

	if category.ID == 0 {
		return apperror.NewUserError("文章分类不存在")
	}

	if err := s.repo.Delete(ctx, category.ID); err != nil {
		return apperror.NewSystemError(err, "删除文章分类失败")
	}

	return nil
}

// ArticleService 文章业务逻辑层
type ArticleService struct {
	repo *repository.ArticleRepository
}

// NewArticleService 创建文章业务逻辑层
func NewArticleService(repo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		repo: repo,
	}
}

// List 获取文章列表
func (s *ArticleService) List(ctx context.Context, req *dto.ArticleListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.Article{},
		Total: 0,
	}

	filter := &repository.ArticleFilter{
		Title:      req.Title,
		CategoryId: req.CategoryId,
		Status:     req.Status,
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文章总数失败")
	}

	result.Total = total
	if total == 0 {
		return result, nil
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "sort", Sort: "asc"},
		{Field: "id", Sort: "desc"},
	}

	// 预加载图片信息
	list, err := s.repo.FindAll(ctx, filter, pagination, orders, repository.WithPreloads("Image"))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询文章列表失败")
	}

	result.List = list
	return result, nil
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, req *dto.ArticleCreateRequest) error {
	article := &model.Article{
		Title:        req.Title,
		ShowType:     req.ShowType,
		Content:      req.Content,
		CategoryId:   req.CategoryId,
		ImageId:      req.ImageId,
		Status:       req.Status,
		Sort:         req.Sort,
		VirtualViews: req.VirtualViews,
		StoreId:      req.StoreId,
	}

	if err := s.repo.Create(ctx, article); err != nil {
		return apperror.NewSystemError(err, "创建文章失败")
	}

	return nil
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, req *dto.ArticleUpdateRequest) error {
	article, err := s.repo.FindOne(ctx, req.Id)
	if err != nil {
		return apperror.NewSystemError(err, "查询文章失败")
	}

	if article.ID == 0 {
		return apperror.NewUserError("文章不存在")
	}

	updateData := map[string]any{
		"title":         req.Title,
		"content":       req.Content,
		"category_id":   req.CategoryId,
		"status":        req.Status,
		"sort":          req.Sort,
		"virtual_views": req.VirtualViews,
		"image_id":      req.ImageId,
	}

	if err := s.repo.Updates(ctx, article, updateData); err != nil {
		return apperror.NewSystemError(err, "更新文章失败")
	}

	return nil
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, req *dto.ArticleDeleteRequest) error {
	article, err := s.repo.FindOne(ctx, req.Id)
	if err != nil {
		return apperror.NewSystemError(err, "查询文章失败")
	}

	if article.ID == 0 {
		return apperror.NewUserError("文章不存在")
	}

	if err := s.repo.Delete(ctx, article.ID); err != nil {
		return apperror.NewSystemError(err, "删除文章失败")
	}

	return nil
}

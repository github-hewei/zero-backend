package article

import (
	"context"
	"errors"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-web/errcode"
)

// CategoryService 文章分类业务逻辑层
type CategoryService struct {
	repo *CategoryRepository
}

// NewCategoryService 创建文章分类业务逻辑层
func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// List 获取文章分类列表
func (s *CategoryService) List(ctx context.Context, req *CategoryListRequest) (*ListResult, error) {
	result := &ListResult{List: []*Category{}, Total: 0}

	filter := &CategoryFilter{StoreId: req.StoreId, Name: req.Name, Status: req.Status}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取文章分类列表失败"))
	}
	result.Total = total
	if total == 0 {
		return result, nil
	}

	orders := baserepo.Orders{{Field: "sort", Sort: "asc"}, {Field: "id", Sort: "desc"}}
	list, err := s.repo.FindAll(ctx, filter, baserepo.NewPagination(req.Page, req.Limit), orders)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取文章分类列表失败"))
	}
	result.List = list
	return result, nil
}

// Create 创建文章分类
func (s *CategoryService) Create(ctx context.Context, req *CategoryCreateRequest) error {
	category := &Category{Name: req.Name, Status: req.Status, Sort: req.Sort, StoreId: req.StoreId}
	if err := s.repo.Create(ctx, category); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建文章分类失败"))
	}
	return nil
}

// Update 更新文章分类
func (s *CategoryService) Update(ctx context.Context, req *CategoryUpdateRequest) error {
	filter := &CategoryFilter{Id: req.Id, StoreId: req.StoreId}
	category, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("文章分类不存在"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新文章分类失败"))
	}
	updateData := map[string]any{"name": req.Name, "status": req.Status, "sort": req.Sort}
	if err := s.repo.Updates(ctx, category, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新文章分类失败"))
	}
	return nil
}

// Delete 删除文章分类
func (s *CategoryService) Delete(ctx context.Context, req *CategoryDeleteRequest) error {
	filter := &CategoryFilter{Id: req.Id, StoreId: req.StoreId}
	category, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("文章分类不存在或无权限访问"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除文章分类失败"))
	}
	if err := s.repo.Delete(ctx, category.ID); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除文章分类失败"))
	}
	return nil
}

// Service 文章业务逻辑层
type Service struct {
	repo *Repository
}

// NewService 创建文章业务逻辑层
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List 获取文章列表
func (s *Service) List(ctx context.Context, req *ListRequest) (*ListResult, error) {
	result := &ListResult{List: []*Article{}, Total: 0}

	filter := &Filter{StoreId: req.StoreId, Title: req.Title, CategoryId: req.CategoryId, Status: req.Status}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取文章列表失败"))
	}
	result.Total = total
	if total == 0 {
		return result, nil
	}

	orders := baserepo.Orders{{Field: "sort", Sort: "asc"}, {Field: "id", Sort: "desc"}}
	list, err := s.repo.FindAll(ctx, filter, baserepo.NewPagination(req.Page, req.Limit), orders)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取文章列表失败"))
	}
	result.List = list
	return result, nil
}

// Create 创建文章
func (s *Service) Create(ctx context.Context, req *CreateRequest) error {
	article := &Article{
		Title: req.Title, ShowType: req.ShowType, Content: req.Content,
		CategoryId: req.CategoryId, ImageId: req.ImageId,
		Status: req.Status, Sort: req.Sort, VirtualViews: req.VirtualViews, StoreId: req.StoreId,
	}
	if err := s.repo.Create(ctx, article); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建文章失败"))
	}
	return nil
}

// Update 更新文章
func (s *Service) Update(ctx context.Context, req *UpdateRequest) error {
	filter := &Filter{Id: req.Id, StoreId: req.StoreId}
	article, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("文章不存在或无权限访问"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新文章失败"))
	}
	updateData := map[string]any{
		"title": req.Title, "content": req.Content, "category_id": req.CategoryId,
		"status": req.Status, "sort": req.Sort, "virtual_views": req.VirtualViews, "image_id": req.ImageId,
	}
	if err := s.repo.Updates(ctx, article, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新文章失败"))
	}
	return nil
}

// Delete 删除文章
func (s *Service) Delete(ctx context.Context, req *DeleteRequest) error {
	filter := &Filter{Id: req.Id, StoreId: req.StoreId}
	article, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("文章不存在或无权限访问"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除文章失败"))
	}
	if err := s.repo.Delete(ctx, article.ID); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除文章失败"))
	}
	return nil
}

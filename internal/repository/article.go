package repository

import (
	"zero-backend/internal/model"
	"zero-backend/pkg/helper"

	"gorm.io/gorm"
)

// ArticleCategoryFilter 文章分类过滤条件
type ArticleCategoryFilter struct {
	Id      uint32
	StoreId uint32
	Name    string
	Status  int8
}

// Apply 应用筛选条件
func (f *ArticleCategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Id != 0 {
		db = db.Where("id = ?", f.Id)
	}

	if f.StoreId != 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	if f.Name != "" {
		db = db.Where("name like ?", f.Name+"%")
	}

	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}

	return db
}

// ArticleCategoryRepository 文章分类数据操作
type ArticleCategoryRepository struct {
	*BaseRepository[model.ArticleCategory]
}

// NewArticleCategoryRepository 创建文章分类数据操作
func NewArticleCategoryRepository(db *gorm.DB) *ArticleCategoryRepository {
	baseRepo := NewBaseRepository[model.ArticleCategory](db)
	return &ArticleCategoryRepository{
		BaseRepository: baseRepo,
	}
}

// ArticleFilter 文章过滤条件
type ArticleFilter struct {
	Id         uint32
	StoreId    uint32
	Title      string
	CategoryId uint32
	Status     int8
}

// Apply 应用筛选条件
func (f *ArticleFilter) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Id != 0 {
		db = db.Where("id = ?", f.Id)
	}

	if f.StoreId != 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	if f.Title != "" {
		db = db.Where("title like ?", helper.SafeLikeString(f.Title)+"%")
	}

	if f.CategoryId != 0 {
		db = db.Where("category_id = ?", f.CategoryId)
	}

	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}

	return db
}

// ArticleRepository 文章数据操作
type ArticleRepository struct {
	*BaseRepository[model.Article]
}

// NewArticleRepository 创建文章数据操作
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	baseRepo := NewBaseRepository[model.Article](db)
	return &ArticleRepository{
		BaseRepository: baseRepo,
	}
}

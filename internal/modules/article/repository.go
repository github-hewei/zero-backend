package article

import (
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"gorm.io/gorm"
)

// CategoryFilter 文章分类过滤条件
type CategoryFilter struct {
	Id      uint32
	StoreId uint32
	Name    string
	Status  int8
}

// Apply 过滤条件
func (f *CategoryFilter) Apply(db *gorm.DB) *gorm.DB {
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

// CategoryRepository 文章分类数据操作
type CategoryRepository struct {
	*baserepo.BaseRepository[Category]
}

// NewCategoryRepository 创建文章分类数据操作
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{BaseRepository: baserepo.NewBaseRepository[Category](db)}
}

// Filter 文章过滤条件
type Filter struct {
	Id         uint32
	StoreId    uint32
	Title      string
	CategoryId uint32
	Status     int8
}

// Apply 过滤条件
func (f *Filter) Apply(db *gorm.DB) *gorm.DB {
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

// Repository 文章数据操作
type Repository struct {
	*baserepo.BaseRepository[Article]
}

// NewRepository 创建文章数据操作
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{BaseRepository: baserepo.NewBaseRepository[Article](db)}
}

package region

import (
	"context"

	"gorm.io/gorm"
)

// Repository 区划数据操作
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建区划数据操作
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindAll 查询所有区划数据
func (r *Repository) FindAll(ctx context.Context) ([]Region, error) {
	var regions []Region
	if err := r.db.WithContext(ctx).Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

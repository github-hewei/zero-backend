package repository

import (
	"context"
	"errors"
	"zero-backend/internal/model"

	"gorm.io/gorm"
)

// RegionRepository 通用数据操作
type RegionRepository struct {
	Db *gorm.DB
}

// NewRegionRepository 创建通用数据操作
func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{Db: db}
}

// FindAll 获取所有地区数据
func (r *RegionRepository) FindAll(ctx context.Context) ([]model.Region, error) {
	var regions []model.Region
	err := r.Db.WithContext(ctx).Find(&regions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return regions, nil
}

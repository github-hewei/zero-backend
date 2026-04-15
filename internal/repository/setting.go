package repository

import (
	"zero-backend/internal/model"

	"gorm.io/gorm"
)

// SettingFilterField 设置表过滤字段
type SettingFilterField struct {
	Id         uint32
	StoreId    uint32
	SettingKey string
}

// Apply 应用筛选条件
func (f *SettingFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Id != 0 {
		db = db.Where("id = ?", f.Id)
	}

	if f.StoreId != 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	if f.SettingKey != "" {
		db = db.Where("setting_key = ?", f.SettingKey)
	}

	return db
}

// SettingRepository 设置数据操作
type SettingRepository struct {
	*BaseRepository[model.Setting]
}

// NewSettingRepository 创建设置数据操作
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	baseRepo := NewBaseRepository[model.Setting](db)
	return &SettingRepository{
		BaseRepository: baseRepo,
	}
}

// SettingDefaultFilterField 默认设置表过滤字段
type SettingDefaultFilterField struct {
	SettingKey string
}

// Apply 应用筛选条件
func (f *SettingDefaultFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.SettingKey != "" {
		db = db.Where("setting_key = ?", f.SettingKey)
	}

	return db
}

// SettingDefaultRepository 默认设置数据操作
type SettingDefaultRepository struct {
	*BaseRepository[model.SettingDefault]
}

// NewSettingDefaultRepository 创建默认设置数据操作
func NewSettingDefaultRepository(db *gorm.DB) *SettingDefaultRepository {
	baseRepo := NewBaseRepository[model.SettingDefault](db)
	return &SettingDefaultRepository{
		BaseRepository: baseRepo,
	}
}

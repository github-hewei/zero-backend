package setting

import (
	"github.com/241x/zero-kit/baserepo"
	"gorm.io/gorm"
)

type Filter struct {
	Id         uint32
	StoreId    uint32
	SettingKey string
}

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
	if f.SettingKey != "" {
		db = db.Where("setting_key = ?", f.SettingKey)
	}
	return db
}

type Repository struct {
	*baserepo.BaseRepository[Setting]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{BaseRepository: baserepo.NewBaseRepository[Setting](db)}
}

type DefaultFilter struct {
	SettingKey string
}

func (f *DefaultFilter) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}
	if f.SettingKey != "" {
		db = db.Where("setting_key = ?", f.SettingKey)
	}
	return db
}

type DefaultRepository struct {
	*baserepo.BaseRepository[DefaultSetting]
}

func NewDefaultRepository(db *gorm.DB) *DefaultRepository {
	return &DefaultRepository{BaseRepository: baserepo.NewBaseRepository[DefaultSetting](db)}
}

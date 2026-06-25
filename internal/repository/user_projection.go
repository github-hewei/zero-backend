package repository

import (
	"zero-backend/internal/model"

	"github.com/241x/zero-kit/baserepo"
	"gorm.io/gorm"
)

// UserFilterField 用户过滤字段（投影）
type UserFilterField struct {
	Id       uint32
	StoreId  uint32
	Username string
}

func (f *UserFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}
	if f.Id != 0 {
		db = db.Where("id = ?", f.Id)
	}
	if f.StoreId != 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}
	if f.Username != "" {
		db = db.Where("username = ?", f.Username)
	}
	return db
}

// UserRepository 用户数据操作（投影）
type UserRepository struct {
	*baserepo.BaseRepository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: baserepo.NewBaseRepository[model.User](db)}
}

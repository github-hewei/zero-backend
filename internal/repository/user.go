package repository

import (
	"time"
	"zero-backend/internal/model"

	"gorm.io/gorm"
)

// UserFilterField 用户表过滤字段
type UserFilterField struct {
	Id       uint32
	StoreId  uint32
	Username string
	Mobile   string
	Status   int8
}

// Apply 实现Filter接口
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

	if f.Mobile != "" {
		db = db.Where("mobile = ?", f.Mobile)
	}

	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}

	return db
}

// UserRepository 用户数据操作
type UserRepository struct {
	*BaseRepository[model.User]
}

// NewUserRepository 创建用户数据操作
func NewUserRepository(db *gorm.DB) *UserRepository {
	baseRepo := NewBaseRepository[model.User](db)
	return &UserRepository{
		BaseRepository: baseRepo,
	}
}

// UserPointsLogRepository 用户积分记录数据操作
type UserPointsLogRepository struct {
	*BaseRepository[model.UserPointsLog]
}

// NewUserPointsLogRepository 创建用户积分记录数据操作
func NewUserPointsLogRepository(db *gorm.DB) *UserPointsLogRepository {
	baseRepo := NewBaseRepository[model.UserPointsLog](db)
	return &UserPointsLogRepository{
		BaseRepository: baseRepo,
	}
}

// UserPointsLogFilterField 积分记录过滤字段
type UserPointsLogFilterField struct {
	StoreId    uint32
	UserId     uint32
	StartDate  string
	EndDate    string
	ChangeType int8
}

// Apply 实现Filter接口
func (f *UserPointsLogFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.StoreId != 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	if f.UserId > 0 {
		db = db.Where("user_id = ?", f.UserId)
	}

	if f.ChangeType > 0 {
		db = db.Where("change_type = ?", f.ChangeType)
	}

	if f.StartDate != "" {
		if startTime, err := time.Parse("2006-01-02", f.StartDate); err == nil {
			db = db.Where("created_at >= ?", startTime.Unix())
		}
	}

	if f.EndDate != "" {
		if endTime, err := time.Parse("2006-01-02", f.EndDate); err == nil {
			db = db.Where("created_at <= ?", endTime.Unix()+86400-1) // 当天23:59:59
		}
	}

	return db
}

package user

import (
	"time"

	"github.com/241x/zero-kit/baserepo"
	"gorm.io/gorm"
)

// Filter 用户表过滤字段
type Filter struct {
	Id       uint32
	StoreId  uint32
	Username string
	Mobile   string
	Status   int8
}

// Apply 应用过滤条件
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

// Repository 用户数据操作
type Repository struct {
	*baserepo.BaseRepository[User]
}

// NewRepository 创建用户数据操作实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{BaseRepository: baserepo.NewBaseRepository[User](db)}
}

// PointsLogFilter 积分记录过滤字段
type PointsLogFilter struct {
	StoreId    uint32
	UserId     uint32
	StartDate  string
	EndDate    string
	ChangeType int8
}

// Apply 应用过滤条件
func (f *PointsLogFilter) Apply(db *gorm.DB) *gorm.DB {
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
			db = db.Where("created_at <= ?", endTime.Unix()+86400-1)
		}
	}
	return db
}

// PointsLogRepo 用户积分记录数据操作
type PointsLogRepo struct {
	*baserepo.BaseRepository[UserPointsLog]
}

// NewPointsLogRepo 创建用户积分记录数据操作实例
func NewPointsLogRepo(db *gorm.DB) *PointsLogRepo {
	return &PointsLogRepo{BaseRepository: baserepo.NewBaseRepository[UserPointsLog](db)}
}

package user

import (
	"context"
	"time"

	"github.com/241x/zero-kit/baserepo"
	"gorm.io/gorm"
)

// RepositoryInterface 用户数据操作接口
type RepositoryInterface interface {
	FindOne(ctx context.Context, id any, opts ...baserepo.QueryOption) (*User, error)
	FindAll(ctx context.Context, filter baserepo.Filter, pagination baserepo.Paginator, orders baserepo.Orders, opts ...baserepo.QueryOption) ([]*User, error)
	Count(ctx context.Context, filter baserepo.Filter, opts ...baserepo.QueryOption) (int64, error)
	Create(ctx context.Context, user *User, opts ...baserepo.CreateOption) error
	Updates(ctx context.Context, user *User, updates map[string]any, opts ...baserepo.UpdateOption) error
	Delete(ctx context.Context, id any, opts ...baserepo.DeleteOption) error
}

// PointsLogRepositoryInterface 用户积分记录数据操作接口
type PointsLogRepositoryInterface interface {
	FindAll(ctx context.Context, filter baserepo.Filter, pagination baserepo.Paginator, orders baserepo.Orders, opts ...baserepo.QueryOption) ([]*UserPointsLog, error)
	Count(ctx context.Context, filter baserepo.Filter, opts ...baserepo.QueryOption) (int64, error)
	Create(ctx context.Context, log *UserPointsLog, opts ...baserepo.CreateOption) error
}

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

// PointsLogRepository 用户积分记录数据操作
type PointsLogRepository struct {
	*baserepo.BaseRepository[UserPointsLog]
}

// NewPointsLogRepository 创建用户积分记录数据操作实例
func NewPointsLogRepository(db *gorm.DB) *PointsLogRepository {
	return &PointsLogRepository{BaseRepository: baserepo.NewBaseRepository[UserPointsLog](db)}
}

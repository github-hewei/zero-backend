package platform_user

import (
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"gorm.io/gorm"
)

// PlatformUserFilterField 平台用户过滤字段
type PlatformUserFilterField struct {
	Username string
	RealName string
}

// Apply 应用过滤条件
func (f *PlatformUserFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}
	if f.Username != "" {
		db = db.Where("username like ?", helper.SafeLikeString(f.Username)+"%")
	}
	if f.RealName != "" {
		db = db.Where("real_name like ?", helper.SafeLikeString(f.RealName)+"%")
	}
	return db
}

// PlatformUserUsernameFilterField 用户名精确查找
type PlatformUserUsernameFilterField struct {
	Username string
}

// Apply 应用过滤条件
func (f *PlatformUserUsernameFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}
	if f.Username != "" {
		db = db.Where("username = ?", f.Username)
	}
	return db
}

// PlatformUserRepository 平台用户数据操作
type PlatformUserRepository struct {
	*baserepo.BaseRepository[PlatformUser]
}

// NewPlatformUserRepository 创建平台用户数据操作
func NewPlatformUserRepository(db *gorm.DB) *PlatformUserRepository {
	return &PlatformUserRepository{BaseRepository: baserepo.NewBaseRepository[PlatformUser](db)}
}

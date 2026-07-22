package platform_user

import (
	"slices"

	"gorm.io/plugin/soft_delete"
)

// PlatformUser 平台管理员表模型
type PlatformUser struct {
	ID            uint32 `json:"id" gorm:"primaryKey"`
	Username      string `json:"username" gorm:"size:64;not null;default:'';comment:用户名;uniqueIndex"`
	Password      string `json:"-" gorm:"size:255;not null;default:'';comment:登录密码"`
	RealName      string `json:"real_name" gorm:"size:64;not null;default:'';comment:姓名"`
	AvatarID      uint32 `json:"avatar_id" gorm:"not null;default:0;comment:头像文件ID"`
	Role          int8   `json:"role" gorm:"type:tinyint;not null;default:0;comment:角色(0超管 1运营 2审计)"`
	Status        int8   `json:"status" gorm:"type:tinyint;not null;default:1;comment:状态(1启用 0禁用)"`
	LastLoginTime int64  `json:"last_login_time" gorm:"not null;default:0;comment:最后登录时间"`
	LastLoginIP   string `json:"last_login_ip" gorm:"size:45;not null;default:'';comment:最后登录IP"`
	CreatedAt     int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt     int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`

	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"not null;default:0;comment:删除时间"`
}

// PlatformRole 平台角色
type PlatformRole int8

const (
	RoleSuperAdmin PlatformRole = 0 // 超级管理员：全部权限
	RoleOperator   PlatformRole = 1 // 运营：租户管理、重置密码、快捷登录
	RoleAuditor    PlatformRole = 2 // 审计：只读权限
)

// HasPermission 检查角色是否拥有指定操作权限
func (r PlatformRole) HasPermission(required ...PlatformRole) bool {
	return slices.Contains(required, r)
}

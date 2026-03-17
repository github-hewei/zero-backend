package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// User 用户表模型
type User struct {
	ID            uint32 `json:"id" gorm:"primaryKey"`
	Username      string `json:"username" gorm:"size:32;not null;default:'';comment:用户名;index:username"`
	Password      string `json:"-" gorm:"size:255;not null;default:'';comment:密码"`
	Mobile        string `json:"mobile" gorm:"size:30;not null;default:'';comment:用户手机号;index:mobile"`
	NickName      string `json:"nick_name" gorm:"size:64;not null;default:'';comment:用户昵称"`
	AvatarId      uint32 `json:"avatar_id" gorm:"not null;default:0;comment:头像文件ID"`
	Gender        int8   `json:"gender" gorm:"type:tinyint;not null;default:0;comment:性别"`
	Country       string `json:"country" gorm:"size:50;not null;default:'';comment:国家"`
	Province      string `json:"province" gorm:"size:50;not null;default:'';comment:省份"`
	City          string `json:"city" gorm:"size:50;not null;default:'';comment:城市"`
	Platform      string `json:"platform" gorm:"size:20;not null;default:'';comment:注册来源"`
	Status        int8   `json:"status" gorm:"type:tinyint;not null;default:1;comment:账号状态: 1正常 2禁用"`
	Points        uint32 `json:"points" gorm:"not null;default:0;comment:用户积分"`
	LastLoginTime uint32 `json:"last_login_time" gorm:"not null;default:0;comment:最后登录时间"`
	StoreId       uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:store_id"`
	CreatedAt     int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt     int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`

	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"not null;default:0;comment:删除时间"`
	Avatar    *UploadFile           `json:"avatar" gorm:"foreignKey:AvatarId"`
}

// TableName 指定数据表名称
func (m *User) TableName() string {
	return TableNamePrefix + "user"
}

// UserPointsLog 用户积分变更记录模型
type UserPointsLog struct {
	ID         uint32 `json:"id" gorm:"primaryKey"`
	StoreId    uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:idx_store_user,priority:1"`
	UserId     uint32 `json:"user_id" gorm:"not null;default:0;comment:用户ID;index:idx_store_user,priority:2"`
	Points     int32  `json:"points" gorm:"not null;default:0;comment:变更积分值"`
	ChangeType int8   `json:"change_type" gorm:"type:tinyint;not null;default:0;comment:变更类型(1增加 2减少)"`
	SourceType int8   `json:"source_type" gorm:"type:tinyint;not null;default:0;comment:来源类型(10消费 20充值 30活动);index:idx_source,priority:1"`
	SourceId   string `json:"source_id" gorm:"size:50;not null;default:'';comment:来源ID(如订单号);index:idx_source,priority:2"`
	Remark     string `json:"remark" gorm:"size:255;not null;default:'';comment:备注"`
	CreatedAt  int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`

	SourceTypeText string `json:"source_type_text" gorm:"-"`
}

// TableName 指定数据表名称
func (m *UserPointsLog) TableName() string {
	return TableNamePrefix + "user_points_log"
}

// AfterFind 查询后处理
func (m *UserPointsLog) AfterFind(tx *gorm.DB) (err error) {
	// m.SourceTypeText = constants.PointsSourceType(m.SourceType).String()
	return nil
}

package model

// User 用户模型投影，仅包含认证模块需要的最小字段集合。
type User struct {
	ID       uint32 `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"-"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	AvatarId uint32 `json:"avatar_id"`
	StoreId  uint32 `json:"store_id"`
	Status   int8   `json:"status"`
}

func (User) TableName() string { return TableNamePrefix + "user" }

package model

// Setting 系统设置
type Setting struct {
	ID            uint32 `json:"id" gorm:"primaryKey"`
	SettingKey    string `json:"setting_key" gorm:"size:30;not null;comment:设置项标识"`
	SettingValues string `json:"setting_values" gorm:"not null;comment:设置内容（json格式）"`
	Description   string `json:"description" gorm:"size:255;not null;default:'';comment:设置项描述"`
	StoreId       uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;uniqueIndex:unique_key"`
	CreatedAt     int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt     int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`
}

// TableName 指定数据表名称
func (m *Setting) TableName() string {
	return TableNamePrefix + "setting"
}

// SettingDefault 系统默认设置
type SettingDefault struct {
	ID            uint32 `json:"id" gorm:"primaryKey"`
	SettingKey    string `json:"setting_key" gorm:"size:30;not null;comment:设置项标识"`
	SettingValues string `json:"setting_values" gorm:"not null;comment:设置内容（json格式）"`
	Description   string `json:"description" gorm:"size:255;not null;default:'';comment:设置项描述"`
	CreatedAt     int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt     int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`
}

// TableName 指定数据表名称
func (m *SettingDefault) TableName() string {
	return TableNamePrefix + "setting_default"
}

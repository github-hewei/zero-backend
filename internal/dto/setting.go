package dto

// SettingListRequest 获取设置列表请求参数
type SettingListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
	StoreId    uint32 `json:"store_id"`
}

// SettingCreateRequest 创建设置请求参数
type SettingCreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

// SettingUpdateRequest 更新设置请求参数
type SettingUpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

// SettingDeleteRequest 删除设置请求参数
type SettingDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// SettingDefaultListRequest 获取默认设置列表请求参数
type SettingDefaultListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
}

// SettingDefaultCreateRequest 创建默认设置请求参数
type SettingDefaultCreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

// SettingDefaultUpdateRequest 更新默认设置请求参数
type SettingDefaultUpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

// SettingDefaultDeleteRequest 删除默认设置请求参数
type SettingDefaultDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// SettingFormConfigsRequest 获取表单配置请求参数
type SettingFormConfigsRequest struct {
	OnlyPlatform bool `json:"only_platform"`
}

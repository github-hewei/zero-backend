package setting

// ListResult 列表数据结构体
type ListResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

// ListRequest 列表请求参数
type ListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
	StoreId    uint32 `json:"store_id"`
}

// CreateRequest 创建请求参数
type CreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

// UpdateRequest 更新请求参数
type UpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

// DeleteRequest 删除请求参数
type DeleteRequest struct {
	ID      uint32 `json:"id" validate:"required"`
	StoreId uint32 `json:"store_id"`
}

// DefaultListRequest 默认列表请求参数
type DefaultListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
}

// DefaultCreateRequest 默认创建请求参数
type DefaultCreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

// DefaultUpdateRequest 默认更新请求参数
type DefaultUpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

// DefaultDeleteRequest 默认删除请求参数
type DefaultDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// FormConfigsRequest 表单配置请求参数
type FormConfigsRequest struct {
	OnlyPlatform bool `json:"only_platform"`
}

package setting

// ListResult 列表数据结构体
type ListResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

type ListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
	StoreId    uint32 `json:"store_id"`
}

type CreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

type UpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
	StoreId       uint32 `json:"store_id"`
}

type DeleteRequest struct {
	ID      uint32 `json:"id" validate:"required"`
	StoreId uint32 `json:"store_id"`
}

type DefaultListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=1000"`
	SettingKey string `json:"setting_key"`
}

type DefaultCreateRequest struct {
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

type DefaultUpdateRequest struct {
	ID            uint32 `json:"id" validate:"required"`
	SettingKey    string `json:"setting_key" validate:"required,min=2,max=30"`
	SettingValues string `json:"setting_values" validate:"required"`
	Description   string `json:"description" validate:"max=255"`
}

type DefaultDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

type FormConfigsRequest struct {
	OnlyPlatform bool `json:"only_platform"`
}

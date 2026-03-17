package dto

// SettingFormGroup 表单组配置
type SettingFormGroup struct {
	Key          string             `json:"key"`           // 设置项key
	Label        string             `json:"label"`         // 显示标签
	Description  string             `json:"description"`   // 描述
	OnlyPlatform bool               `json:"only_platform"` // 是否只对平台可见
	Fields       []SettingFormField `json:"fields"`        // 表单字段
}

// SettingFormField 表单字段配置
type SettingFormField struct {
	Key      string       `json:"key"`      // 字段key
	Label    string       `json:"label"`    // 显示标签
	Type     string       `json:"type"`     // 表单类型
	Required bool         `json:"required"` // 是否必填
	Options  []FormOption `json:"options"`  // 选项
}

// FormOption 表单选项
type FormOption struct {
	Label string      `json:"label"` // 显示文本
	Value interface{} `json:"value"` // 选项值
}

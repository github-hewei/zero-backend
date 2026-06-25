package setting

// FormGroup 表单组配置
type FormGroup struct {
	Key          string      `json:"key"`
	Label        string      `json:"label"`
	Description  string      `json:"description"`
	OnlyPlatform bool        `json:"only_platform"`
	Fields       []FormField `json:"fields"`
}

// FormField 表单字段配置
type FormField struct {
	Key      string       `json:"key"`
	Label    string       `json:"label"`
	Type     string       `json:"type"`
	Required bool         `json:"required"`
	Options  []FormOption `json:"options"`
}

// FormOption 表单选项
type FormOption struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// QiniuTokenResponse 七牛上传token
type QiniuTokenResponse struct {
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	UploadUrl string `json:"upload_url"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
	Zone      string `json:"zone"`
	IsEnabled bool   `json:"is_enabled"`
}

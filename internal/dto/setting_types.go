package dto

// SiteConfig 站点信息配置
type SiteConfig struct {
	SiteName    string `json:"site_name"`
	Logo        string `json:"logo"`
	Favicon     string `json:"favicon"`
	Copyright   string `json:"copyright"`
	ICP         string `json:"icp"`
	Maintenance bool   `json:"maintenance"`
}

// EmailConfig 邮件服务配置
type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	SSL      bool   `json:"ssl"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
	IsEnabled bool   `json:"is_enabled"`
	Zone      string `json:"zone"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	StorageType  string   `json:"storage_type"`  // 存储方式(local/qiniu)
	MaxSize      string   `json:"max_size"`      // 最大文件大小(MB)
	AllowedTypes []string `json:"allowed_types"` // 允许的文件类型
	PathFormat   string   `json:"path_format"`   // 存储路径格式
	IsEnabled    bool     `json:"is_enabled"`    // 是否启用
}

// DemoConfig 测试配置(包含所有字段类型)
type DemoConfig struct {
	// 基础字段
	TextInput   string `json:"text_input"`   // 文本输入
	TextArea    string `json:"text_area"`    // 多行文本
	NumberInput int    `json:"number_input"` // 数字输入

	// 选择类字段
	SelectField   string   `json:"select_field"`   // 下拉选择
	RadioField    string   `json:"radio_field"`    // 单选按钮
	CheckboxField []string `json:"checkbox_field"` // 多选框

	// 开关类
	SwitchField bool `json:"switch_field"` // 开关

	// 复杂类型
	ObjectField struct {
		SubField1 string `json:"sub_field1"`
		SubField2 int    `json:"sub_field2"`
	} `json:"object_field"`

	ArrayField []struct {
		ItemName  string `json:"item_name"`
		ItemValue int    `json:"item_value"`
	} `json:"array_field"`

	// 日期时间
	DateField     string `json:"date_field"`     // 日期选择
	DateTimeField string `json:"datetime_field"` // 日期时间选择
}

package captcha

// GenerateResponse 验证码生成响应
type GenerateResponse struct {
	CaptchaID   string `json:"captcha_id"`
	MasterImage string `json:"master_image"`
	ThumbImage  string `json:"thumb_image"`
}

// ClickPoint 验证码点选坐标（前端提交）
type ClickPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Config 验证码配置
type Config struct {
	Enabled bool
	TTL     int
}

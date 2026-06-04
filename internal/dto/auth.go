package dto

import "zero-backend/internal/model"

// AuthLoginRequest 登录参数
type AuthLoginRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CaptchaID   string `json:"captcha_id" validate:"required"`
	CaptchaCode string `json:"captcha_code" validate:"required"`
}

// AuthGetPermissionsRequest 获取权限参数
type AuthGetPermissionsRequest struct {
	IsTree bool `json:"is_tree"`
}

// ChangePasswordRequest 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=20"`
}

// CaptchaGenerateResponse 验证码生成响应（返回给前端）
type CaptchaGenerateResponse struct {
	CaptchaID   string `json:"captcha_id"`   // 验证码唯一标识
	MasterImage string `json:"master_image"` // 主图 base64（含文字提示的大图）
	ThumbImage  string `json:"thumb_image"`  // 缩略图 base64（需要点击的小图）
}

// CaptchaClickPoint 验证码点选坐标（前端提交）
type CaptchaClickPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string      `json:"token"`
	Ttl   int         `json:"ttl"`
	User  *model.User `json:"user,omitempty"`
}

// AdminLoginResponse 登录返回用户信息
type AdminLoginResponse struct {
	Token string          `json:"token"`
	Ttl   int             `json:"ttl"`
	User  *model.RbacUser `json:"user,omitempty"`
}

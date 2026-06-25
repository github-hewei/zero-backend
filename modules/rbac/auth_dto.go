package rbac

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

// AdminLoginResponse 登录返回
type AdminLoginResponse struct {
	Token string `json:"token"`
	Ttl   int    `json:"ttl"`
	User  *RbacUser `json:"user,omitempty"`
}

// CaptchaGenerateResponse 验证码生成响应
type CaptchaGenerateResponse struct {
	CaptchaID   string `json:"captcha_id"`
	MasterImage string `json:"master_image"`
	ThumbImage  string `json:"thumb_image"`
}

// CaptchaClickPoint 验证码点选坐标
type CaptchaClickPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

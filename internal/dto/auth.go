package dto

import "zero-backend/internal/model"

// AuthLoginRequest 登录参数
type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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

package dto

import "zero-backend/internal/model"

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string      `json:"token"`
	Ttl   int         `json:"ttl"`
	User  *model.User `json:"user,omitempty"`
}

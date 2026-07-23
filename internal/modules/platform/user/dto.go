package user

// LoginRequest 平台登录请求
type LoginRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CaptchaID   string `json:"captcha_id" validate:"required"`
	CaptchaCode string `json:"captcha_code" validate:"required"`
}

// LoginResponse 平台登录响应
type LoginResponse struct {
	Token string        `json:"token"`
	Ttl   int           `json:"ttl"`
	User  *PlatformUser `json:"user,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=20"`
}

// UserListRequest 平台用户列表请求
type UserListRequest struct {
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Page     int    `json:"page" validate:"required,min=1"`
	Limit    int    `json:"limit" validate:"required,min=1,max=100"`
}

// UserCreateRequest 创建平台用户请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=5,max=64,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	RealName string `json:"real_name" validate:"required,min=2,max=64"`
	Role     int8   `json:"role" validate:"oneof=0 1 2"`
	AvatarID uint32 `json:"avatar_id"`
}

// UserUpdateRequest 更新平台用户请求
type UserUpdateRequest struct {
	ID       uint32 `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=5,max=64,alphanum"`
	RealName string `json:"real_name" validate:"required,min=2,max=64"`
	Password string `json:"password" validate:"max=20"`
	Role     int8   `json:"role" validate:"oneof=0 1 2"`
	Status   int8   `json:"status" validate:"oneof=0 1"`
	AvatarID uint32 `json:"avatar_id"`
}

// UserDeleteRequest 删除平台用户请求
type UserDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// UserResetPasswordRequest 重置平台用户密码请求
type UserResetPasswordRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// ListResult 列表数据响应
type ListResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

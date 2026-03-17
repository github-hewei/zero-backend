package apperror

// NewError 创建通用错误
func NewError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewUserError 创建用户错误
func NewUserError(message string) *UserError {
	return &UserError{
		AppError: AppError{
			Code:    ErrorCodeUser,
			Message: message,
			Err:     nil,
		},
	}
}

// NewSystemError 创建系统错误
func NewSystemError(err error, message string) *SystemError {
	return &SystemError{
		AppError: AppError{
			Code:    ErrorCodeSystem,
			Message: message,
			Err:     err,
		},
	}
}

// NewUnauthorizedError 创建权限错误
func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{
		AppError: AppError{
			Code:    ErrorCodeUnauthorized,
			Message: "资源未授权，请先登录",
		},
	}
}

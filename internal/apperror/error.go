package apperror

type ErrorCode int

const (
	ErrorCodeNone         ErrorCode = 0
	ErrorCodeUser         ErrorCode = 4000
	ErrorCodeUnauthorized ErrorCode = 4001
	ErrorCodeSystem       ErrorCode = 5000
)

// AppError 错误
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
	TraceId string
}

// Error 错误接口
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 获取内部错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithCode 设置新的错误码
func (e *AppError) WithCode(code ErrorCode) *AppError {
	e.Code = code
	return e
}

// WithMessage 设置新的错误消息
func (e *AppError) WithMessage(msg string) *AppError {
	e.Message = msg
	return e
}

// UserError 用户级别错误
type UserError struct {
	AppError
}

// SystemError 系统级别错误
type SystemError struct {
	AppError
}

// UnauthorizedError 权限错误
type UnauthorizedError struct {
	AppError
}

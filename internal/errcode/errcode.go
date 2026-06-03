package errcode

import "zero-backend/pkg/apperror"

var (
	// InvalidInput 输入校验失败
	InvalidInput = apperror.NewCode(400001, "INVALID_INPUT", "请求参数有误")

	// Unauthorized 未认证
	Unauthorized = apperror.NewCode(401001, "UNAUTHORIZED", "请先登录")

	// Forbidden 无权限
	Forbidden = apperror.NewCode(403001, "FORBIDDEN", "无权限执行此操作")

	// NotFound 资源不存在
	NotFound = apperror.NewCode(404001, "NOT_FOUND", "找不到此记录")

	// Conflict 资源冲突
	Conflict = apperror.NewCode(409001, "CONFLICT", "数据已存在")

	// Internal 系统错误
	Internal = apperror.NewCode(500001, "INTERNAL", "系统异常，请稍后重试")
)

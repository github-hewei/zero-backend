package errcode

import "zero-backend/pkg/apperror"

var (
	// Unauthorized 未认证
	Unauthorized = apperror.NewCode(4001, "UNAUTHORIZED", "您还未登录，请先登录")

	// InvalidInput 输入校验失败
	InvalidInput = apperror.NewCode(4002, "INVALID_INPUT", "输入参数不正确")

	// Forbidden 无权限
	Forbidden = apperror.NewCode(4003, "FORBIDDEN", "无权限执行此操作")

	// NotFound 资源不存在
	NotFound = apperror.NewCode(4004, "NOT_FOUND", "找不到此记录")

	// Conflict 资源冲突
	Conflict = apperror.NewCode(4009, "CONFLICT", "资源已存在")

	// Internal 系统错误
	Internal = apperror.NewCode(5000, "INTERNAL", "系统异常，请稍后重试")
)

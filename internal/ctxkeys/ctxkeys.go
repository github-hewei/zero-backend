package ctxkeys

import (
	"context"
	"time"
)

// traceIDKey 上下文传递请求链路ID
type traceIDKey struct{}

// WithTraceID 注入 traceID
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, id)
}

// TraceID 读取 traceID
func TraceID(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey{}).(string); ok {
		return v
	}
	return ""
}

// beginTimeKey 上下文传递请求开始时间
type beginTimeKey struct{}

// WithBeginTime 注入请求开始时间
func WithBeginTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, beginTimeKey{}, t)
}

// BeginTime 读取请求开始时间
func BeginTime(ctx context.Context) (time.Time, bool) {
	v, ok := ctx.Value(beginTimeKey{}).(time.Time)
	return v, ok
}

// userKey 上下文传递用户信息
type userKey struct{}

// WithUser 注入用户信息
func WithUser(ctx context.Context, user any) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// User 读取用户信息，调用方需自行类型断言
func User(ctx context.Context) any {
	return ctx.Value(userKey{})
}

// storeIdKey 上下文传递企业ID
type storeIdKey struct{}

// WithStoreID 注入企业ID
func WithStoreID(ctx context.Context, id uint32) context.Context {
	return context.WithValue(ctx, storeIdKey{}, id)
}

// StoreID 读取企业ID
func StoreID(ctx context.Context) uint32 {
	v, ok := ctx.Value(storeIdKey{}).(uint32)
	if !ok {
		return 0
	}
	return v
}

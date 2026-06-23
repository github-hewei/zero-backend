package ctxkeys

import (
	"context"
	"time"

	"zero-backend/internal/model"

	webctx "github.com/241x/zero-web/ctxkeys"
)

// 委托到 zero-web/ctxkeys
func WithTraceID(ctx context.Context, id string) context.Context { return webctx.WithTraceID(ctx, id) }
func TraceID(ctx context.Context) string                         { return webctx.TraceID(ctx) }
func WithBeginTime(ctx context.Context, t time.Time) context.Context { return webctx.WithBeginTime(ctx, t) }
func BeginTime(ctx context.Context) (time.Time, bool)               { return webctx.BeginTime(ctx) }
func WithUser(ctx context.Context, user any) context.Context        { return webctx.WithUser(ctx, user) }
func User(ctx context.Context) any                                  { return webctx.User(ctx) }

// UserID 从上下文中获取用户 ID，兼容 RbacUser 与 User 两种模型。
func UserID(ctx context.Context) uint32 {
	if user, ok := User(ctx).(*model.RbacUser); ok {
		return user.ID
	}
	if user, ok := User(ctx).(*model.User); ok {
		return user.ID
	}
	return 0
}

// IsSuperUser 判断当前用户是否为超级管理员。
func IsSuperUser(ctx context.Context) bool {
	if user, ok := User(ctx).(*model.RbacUser); ok {
		return user.SU
	}
	return false
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

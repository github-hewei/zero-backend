package ctxkeys

// TraceIDKey 上下文传递请求链路ID
type TraceIDKey struct{}

// BeginTimeKey 上下文传递请求开始时间
type BeginTimeKey struct{}

// UserKey 上下文传递用户信息
type UserKey struct{}

// QiniuConfigKey 上下文传递七牛云配置信息
type QiniuConfigKey struct{}

// StoreIdKey 上下文传递企业ID
type StoreIdKey struct{}

// CLIContextKey CLI 上下文键
type CLIContextKey struct{}

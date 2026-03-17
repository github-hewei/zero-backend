package logger

import "context"

// nopLogger 空操作 Logger
type nopLogger struct{}

// Nop 返回一个空操作的 Logger
func Nop() Logger {
	return &nopLogger{}
}

// Debug 记录 Debug 级别日志空实现
func (n *nopLogger) Debug(msg string, fields ...any) {}

// Info 记录 Info 级别日志空实现
func (n *nopLogger) Info(msg string, fields ...any) {}

// Warn 记录 Warn 级别日志空实现
func (n *nopLogger) Warn(msg string, fields ...any) {}

// Error 记录 Error 级别日志空实现
func (n *nopLogger) Error(msg string, fields ...any) {}

// Err 错误日志空实现
func (n *nopLogger) Err(err error, msg string, fields ...any) {}

// Log 统一日志记录空实现
func (n *nopLogger) Log(level Level, msg string, fields ...any) {}

// WithContext 返回原始上下文
func (n *nopLogger) WithContext(ctx context.Context) context.Context { return ctx }

// With 创建一个新的 Logger
func (n *nopLogger) With(fields ...any) Logger { return n }

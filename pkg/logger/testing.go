package logger

import "context"

// MockLogger 实现 Logger 接口，用于单元测试
type MockLogger struct {
	Logs []MockLogEntry
}

// MockLogEntry 记录单条日志
type MockLogEntry struct {
	Level  Level
	Msg    string
	Fields []any
}

// NewMockLogger 创建 MockLogger 实例
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Logs: make([]MockLogEntry, 0),
	}
}

// Debug 记录 Debug 级别日志
func (m *MockLogger) Debug(msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{DebugLevel, msg, fields})
}

// Info 记录 Info 级别日志
func (m *MockLogger) Info(msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{InfoLevel, msg, fields})
}

// Warn 记录 Warn 级别日志
func (m *MockLogger) Warn(msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{WarnLevel, msg, fields})
}

// Error 记录 Error 级别日志
func (m *MockLogger) Error(msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{ErrorLevel, msg, fields})
}

// Err 记录包含错误的日志
func (m *MockLogger) Err(err error, msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{ErrorLevel, msg, append(fields, "error", err)})
}

// Log 指定级别记录日志
func (m *MockLogger) Log(level Level, msg string, fields ...any) {
	m.Logs = append(m.Logs, MockLogEntry{level, msg, fields})
}

// WithContext 返回原始 context（简化实现）
func (m *MockLogger) WithContext(ctx context.Context) context.Context {
	return ctx
}

// With 返回自身（简化实现）
func (m *MockLogger) With(fields ...any) Logger {
	return m
}

// HasLog 检查是否存在指定级别的日志
func (m *MockLogger) HasLog(level Level) bool {
	for _, log := range m.Logs {
		if log.Level == level {
			return true
		}
	}
	return false
}

// HasLogMsg 检查是否存在指定级别和消息的日志
func (m *MockLogger) HasLogMsg(level Level, msgContains string) bool {
	for _, log := range m.Logs {
		if log.Level == level && contains(log.Msg, msgContains) {
			return true
		}
	}
	return false
}

// CountByLevel 统计指定级别的日志数量
func (m *MockLogger) CountByLevel(level Level) int {
	count := 0
	for _, log := range m.Logs {
		if log.Level == level {
			count++
		}
	}
	return count
}

// Reset 清空所有日志记录
func (m *MockLogger) Reset() {
	m.Logs = make([]MockLogEntry, 0)
}

// LastLog 获取最后一条日志
func (m *MockLogger) LastLog() *MockLogEntry {
	if len(m.Logs) == 0 {
		return nil
	}
	return &m.Logs[len(m.Logs)-1]
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

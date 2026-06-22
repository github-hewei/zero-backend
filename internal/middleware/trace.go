package middleware

import (
	"time"

	"zero-backend/internal/ctxkeys"

	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceMiddleware 请求链路追踪中间件，注入 traceID、beginTime、logger 到上下文。
type TraceMiddleware struct {
	logger logger.Logger
}

// NewTraceMiddleware 创建 TraceMiddleware
func NewTraceMiddleware(logger logger.Logger) *TraceMiddleware {
	return &TraceMiddleware{logger: logger}
}

// Handle 处理请求
func (m *TraceMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		ctx := c.Request.Context()
		ctx = ctxkeys.WithTraceID(ctx, traceID)
		ctx = gormutil.WithTraceID(ctx, traceID)
		ctx = ctxkeys.WithBeginTime(ctx, time.Now())

		l := m.logger.With("traceId", traceID)
		ctx = l.WithContext(ctx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

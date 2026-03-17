package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"zero-backend/internal/ctxkeys"
	"zero-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BeforeMiddleware struct {
	logger logger.Logger
}

func NewBeforeMiddleware(logger logger.Logger) *BeforeMiddleware {
	return &BeforeMiddleware{
		logger: logger,
	}
}

// Handle 请求前置处理
func (m *BeforeMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := uuid.New().String()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkeys.TraceIDKey{}, traceId)
		ctx = context.WithValue(ctx, ctxkeys.BeginTimeKey{}, time.Now())

		// 添加日志上下文
		logger := m.logger.With("traceId", traceId)
		ctx = logger.WithContext(ctx)
		c.Request = c.Request.WithContext(ctx)

		buffer, err := io.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// 解析请求参数
		postData := map[string]any{}
		if c.Request.Method == http.MethodPost && c.Request.Header.Get("Content-Type") == "application/json" {
			_ = json.Unmarshal(buffer, &postData)
		}

		// 记录请求日志
		logger.Info("Request",
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"query", c.Request.URL.RawQuery,
			"postData", postData,
			"ip", c.ClientIP(),
			"userAgent", c.Request.UserAgent(),
		)

		c.Request.Body = io.NopCloser(bytes.NewReader(buffer))
		c.Next()
	}
}

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/241x/zero-kit/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件，记录请求详情。
type RequestLogger struct{}

// NewRequestLogger 创建请求日志中间件
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Handle 处理请求
func (*RequestLogger) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Ctx(c.Request.Context()).Info("Request",
			"url", c.Request.URL.Path,
			"method", c.Request.Method,
			"query", c.Request.URL.RawQuery,
			"postData", readBody(c),
			"ip", c.ClientIP(),
			"userAgent", c.Request.UserAgent(),
		)

		c.Next()
	}
}

// readBody 读取请求体
func readBody(c *gin.Context) map[string]any {
	ct := strings.ToLower(c.Request.Header.Get("Content-Type"))
	if c.Request.Method != http.MethodPost || !strings.HasPrefix(ct, "application/json") {
		return nil
	}

	buffer, err := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(buffer))
	if err != nil {
		logger.Ctx(c.Request.Context()).Warn("read request body failed", "err", err)
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(buffer, &data); err != nil {
		logger.Ctx(c.Request.Context()).Warn("parse request body failed", "err", err)
		return nil
	}
	return data
}

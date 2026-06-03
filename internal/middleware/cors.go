package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"zero-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
	cfg config.AdminCorsConfig
}

func NewCorsMiddleware(cfg config.AdminCorsConfig) *CorsMiddleware {
	return &CorsMiddleware{
		cfg: cfg,
	}
}

// Handle 跨域规则
func (m *CorsMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		if len(m.cfg.AllowOrigins) > 0 {
			c.Header("Access-Control-Allow-Origin", strings.Join(m.cfg.AllowOrigins, ", "))
		}
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(m.cfg.AllowCredentials))
		c.Header("Access-Control-Allow-Headers", strings.Join(m.cfg.AllowHeaders, ", "))
		c.Header("Access-Control-Allow-Methods", strings.Join(m.cfg.AllowMethods, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

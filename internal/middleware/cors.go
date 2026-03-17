package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"zero-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
	cfg *config.Config
}

func NewCorsMiddleware(config *config.Config) *CorsMiddleware {
	return &CorsMiddleware{
		cfg: config,
	}
}

// Handle 跨域规则
func (m *CorsMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		if len(m.cfg.Admin.Cors.AllowOrigins) > 0 {
			c.Header("Access-Control-Allow-Origin", strings.Join(m.cfg.Admin.Cors.AllowOrigins, ", "))
		}
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(m.cfg.Admin.Cors.AllowCredentials))
		c.Header("Access-Control-Allow-Headers", strings.Join(m.cfg.Admin.Cors.AllowHeaders, ", "))
		c.Header("Access-Control-Allow-Methods", strings.Join(m.cfg.Admin.Cors.AllowMethods, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

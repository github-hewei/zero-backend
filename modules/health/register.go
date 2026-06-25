package health

import (
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// Register 注册健康检查路由
func Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "ok", gin.H{"status": "UP"})
	})
}

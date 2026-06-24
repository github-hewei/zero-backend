package middleware

import (
	"zero-backend/internal/config"
	"zero-backend/modules/api/service"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 权限验证中间件
type AuthMiddleware struct {
	config   config.ApiAuthConfig
	authServ *service.AuthService
}

// NewAuthMiddleware 创建权限验证中间件
func NewAuthMiddleware(cfg config.ApiAuthConfig, authServ *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		config:   cfg,
		authServ: authServ,
	}
}

// LoadUser 从 JWT claims 加载用户信息并注入上下文。
func (m *AuthMiddleware) LoadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := middleware.GetJWTUserID(c)
		ctx := c.Request.Context()

		user, err := m.authServ.GetUserInfo(ctx, userId)
		if err != nil {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		ctx = ctxkeys.WithUser(ctx, user)
		ctx = ctxkeys.WithStoreID(ctx, user.StoreId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

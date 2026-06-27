package platform_user

import (
	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 平台权限验证中间件
type AuthMiddleware struct {
	config   Config
	authServ *AuthService
}

// NewAuthMiddleware 创建平台权限验证中间件
func NewAuthMiddleware(cfg Config, authServ *AuthService) *AuthMiddleware {
	return &AuthMiddleware{config: cfg, authServ: authServ}
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

		if user.Status != 1 {
			response.Error(c, apperror.New(errcode.Forbidden, apperror.WithMsg("账号已被禁用")))
			c.Abort()
			return
		}

		ctx = ctxkeys.WithUser(ctx, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireRole 检查平台用户角色权限
func (m *AuthMiddleware) RequireRole(roles ...PlatformRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := ctxkeys.User(c.Request.Context()).(*PlatformUser)
		if !ok || user == nil {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		if !PlatformRole(user.Role).HasPermission(roles...) {
			response.Error(c, apperror.New(errcode.Forbidden, apperror.WithMsg("暂无权限")))
			c.Abort()
			return
		}

		c.Next()
	}
}

package rbac

import (
	"strings"

	"zero-backend/config"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 权限验证中间件
type AuthMiddleware struct {
	config   config.AdminAuthConfig
	authServ *AuthService
}

// NewAuthMiddleware 创建权限验证中间件
func NewAuthMiddleware(cfg config.AdminAuthConfig, authServ *AuthService) *AuthMiddleware {
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

		ctx = ctxkeys.WithUser(ctx, user)
		ctx = ctxkeys.WithStoreID(ctx, user.StoreId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// CheckAPIPermission 验证接口权限
func (m *AuthMiddleware) CheckAPIPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		user, ok := ctxkeys.User(ctx).(*RbacUser)
		if !ok || user == nil {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		if user.SU {
			c.Next()
			return
		}

		apiPath := strings.TrimPrefix(c.Request.URL.Path, "/api")
		hasPerm, err := m.authServ.CheckAPIPermission(ctx, user, apiPath)
		if err != nil || !hasPerm {
			response.Error(c, apperror.New(errcode.Forbidden, apperror.WithMsg("暂无访问权限，请联系管理员授权")))
			c.Abort()
			return
		}

		c.Next()
	}
}

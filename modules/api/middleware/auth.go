package middleware

import (
	"zero-backend/internal/config"
	"zero-backend/internal/ctxkeys"
	"zero-backend/modules/api/service"

	"github.com/241x/zero-kit/apperror"
	basectxkeys "github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/241x/zero-web/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// JWTAuth 验证JWT
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" || len(tokenString) < 10 {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (any, error) {
			return []byte(m.config.HmacSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			// zlog.Err(err).Str("token", tokenString).Msg("parse jwt error")
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		userId, ok := claims["user_id"]
		if !ok {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		user, err := m.authServ.GetUserInfo(ctx, uint32(userId.(float64)))
		if err != nil {
			response.Error(c, apperror.New(errcode.Unauthorized))
			c.Abort()
			return
		}

		ctx = basectxkeys.WithUser(ctx, user)
		ctx = ctxkeys.WithStoreID(ctx, user.StoreId)
		c.Request = c.Request.WithContext(ctx)
	}
}

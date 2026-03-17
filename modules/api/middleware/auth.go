package middleware

import (
	"context"
	"zero-backend/internal/apperror"
	"zero-backend/internal/config"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/response"
	"zero-backend/modules/api/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 权限验证中间件
type AuthMiddleware struct {
	config   *config.Config
	authServ *service.AuthService
}

// NewAuthMiddleware 创建权限验证中间件
func NewAuthMiddleware(cfg *config.Config, authServ *service.AuthService) *AuthMiddleware {
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
			response.Error(c, apperror.NewUnauthorizedError())
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			return []byte(m.config.Api.HmacSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			// zlog.Err(err).Str("token", tokenString).Msg("parse jwt error")
			response.Error(c, apperror.NewUnauthorizedError())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, apperror.NewUnauthorizedError())
			c.Abort()
			return
		}

		userId, ok := claims["user_id"]
		if !ok {
			response.Error(c, apperror.NewUnauthorizedError())
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		user, err := m.authServ.GetUserInfo(ctx, uint32(userId.(float64)))
		if err != nil {
			response.Error(c, apperror.NewUnauthorizedError())
			c.Abort()
			return
		}

		ctx = context.WithValue(ctx, ctxkeys.UserKey{}, user)
		ctx = context.WithValue(ctx, ctxkeys.StoreIdKey{}, user.StoreId)
		c.Request = c.Request.WithContext(ctx)
	}
}

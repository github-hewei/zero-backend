package platform

import (
	"zero-backend/internal/modules/captcha"
	"zero-backend/internal/modules/platform_user"
	"zero-backend/internal/modules/rbac"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/provider"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewGin 创建一个 gin.Engine 实例
func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(provider.LoadPlatformCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	public := r.Group("/api")
	protected := public.Group("")

	captchaSvc, err := provider.NewCaptchaService(rdb, provider.LoadPlatformCaptchaConfig())
	if err != nil {
		panic(err)
	}
	captcha.Register(public, binder, captchaSvc)

	authCfg := platform_user.MustLoadConfig()
	authMid := platform_user.Register(public, protected, platform_user.Deps{
		DB:      db,
		Binder:  binder,
		Config:  authCfg,
		RDB:     rdb,
		Captcha: captchaSvc,
	})

	platformGroup := protected.Group("")
	platformGroup.Use(authMid.RequireRole(platform_user.RoleSuperAdmin, platform_user.RoleOperator))

	rbac.RegisterPlatform(platformGroup, rbac.PlatformDeps{
		DB:     db,
		Binder: binder,
		Config: rbac.Config{
			HmacSecret:      authCfg.HmacSecret,
			AccessTokenTtl:  authCfg.AccessTokenTtl,
			RefreshTokenTtl: authCfg.RefreshTokenTtl,
		},
		RDB: rdb,
	})

	setting.RegisterPlatform(platformGroup, setting.Deps{DB: db, Binder: binder})

	return r
}

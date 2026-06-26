package router

import (
	"zero-backend/internal/config"
	"zero-backend/internal/modules/region"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/modules/upload"
	"zero-backend/internal/modules/user"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	basecfg "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
	authConfig config.ApiAuthConfig,
	corsConfig basecfg.CorsConfig,
	settingSvc *setting.Service,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	authMid, h := user.RegisterApi(apiGroup, user.ApiDeps{
		DB:      db,
		Binder:  binder,
		AuthCfg: authConfig,
		RDB:     rdb,
	})

	apiGroup.Use(middleware.JWTGuard(authConfig.HmacSecret))
	apiGroup.Use(authMid.LoadUser())

	user.RegisterApiProtected(apiGroup, h)

	upload.RegisterApi(apiGroup, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	setting.RegisterApi(apiGroup, setting.Deps{DB: db, Binder: binder})
	region.Register(apiGroup, region.Deps{DB: db, Binder: binder})

	return r
}

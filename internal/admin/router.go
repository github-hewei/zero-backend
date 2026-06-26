package admin

import (
	"net/http"
	"zero-backend/internal/modules/article"
	"zero-backend/internal/modules/captcha"
	"zero-backend/internal/modules/health"
	"zero-backend/internal/modules/rbac"
	"zero-backend/internal/modules/region"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/modules/upload"
	"zero-backend/internal/modules/user"

	"zero-backend/internal/config"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	basecfg "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	authConfig config.AdminAuthConfig,
	corsConfig basecfg.CorsConfig,
	settingSvc *setting.Service,
	captchaSvc *captcha.Service,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	authMid, h := rbac.RegisterAdmin(apiGroup, rbac.Deps{
		DB:      db,
		Binder:  binder,
		AuthCfg: authConfig,
		RDB:     captchaSvc.RDB(),
		Captcha: captchaSvc,
	})

	captcha.RegisterWith(apiGroup, binder, captchaSvc)

	apiGroup.Use(middleware.JWTGuard(authConfig.HmacSecret))
	apiGroup.Use(authMid.LoadUser())
	apiGroup.Use(authMid.CheckAPIPermission())

	rbac.RegisterAdminProtected(apiGroup, h)

	setting.RegisterAdmin(apiGroup, setting.Deps{DB: db, Binder: binder})
	article.Register(apiGroup, article.Deps{DB: db, Binder: binder})
	upload.RegisterAdmin(apiGroup, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	user.Register(apiGroup, user.Deps{DB: db, Binder: binder})
	region.Register(apiGroup, region.Deps{DB: db, Binder: binder})

	health.Register(r)

	r.LoadHTMLGlob("./views/*.html")
	r.Static("/assets", "./views/assets")
	r.Static("/uploads", "./uploads")

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./views/favicon.ico")
	})
	r.GET("/logo.svg", func(c *gin.Context) {
		c.File("./views/logo.svg")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r
}

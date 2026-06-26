package admin

import (
	"net/http"
	"zero-backend/internal/app"
	"zero-backend/internal/modules/article"
	"zero-backend/internal/modules/captcha"
	"zero-backend/internal/modules/health"
	"zero-backend/internal/modules/rbac"
	"zero-backend/internal/modules/region"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/modules/upload"
	"zero-backend/internal/modules/user"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	captchaSvc *captcha.Service,
	rdb *redis.Client,
) (*gin.Engine, error) {
	r := gin.Default()
	r.Use(middleware.CORS(app.LoadAdminCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	authCfg, err := rbac.LoadConfig()
	if err != nil {
		return nil, err
	}
	authMid, h := rbac.RegisterAdmin(apiGroup, rbac.Deps{
		DB:      db,
		Binder:  binder,
		Config:  authCfg,
		RDB:     captchaSvc.RDB(),
		Captcha: captchaSvc,
	})

	captcha.RegisterWith(apiGroup, binder, captchaSvc)

	apiGroup.Use(middleware.JWTGuard(authCfg.HmacSecret))
	apiGroup.Use(authMid.LoadUser())
	apiGroup.Use(authMid.CheckAPIPermission())

	rbac.RegisterAdminProtected(apiGroup, h)

	settingSvc := app.NewSettingService(db)

	setting.RegisterAdmin(apiGroup, setting.Deps{DB: db, Binder: binder})
	article.Register(apiGroup, article.Deps{DB: db, Binder: binder})
	upload.RegisterAdmin(apiGroup, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	user.Register(apiGroup, user.Deps{DB: db, Binder: binder})
	region.Register(apiGroup, region.Deps{DB: db, Binder: binder})

	health.Register(r)

	r.LoadHTMLGlob("./views/*.html")
	r.Static("/assets", "./views/assets")
	r.Static("/uploads", "./uploads")

	r.GET("/favicon.ico", func(c *gin.Context) { c.File("./views/favicon.ico") })
	r.GET("/logo.svg", func(c *gin.Context) { c.File("./views/logo.svg") })
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.NoRoute(func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })

	return r, nil
}

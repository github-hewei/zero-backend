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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(app.LoadAdminCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	public := r.Group("/api")
	protected := public.Group("")

	captchaSvc := app.Must(app.NewCaptchaService(rdb, app.LoadCaptchaConfig()))
	captcha.RegisterWith(public, binder, captchaSvc)

	authCfg := app.Must(rbac.LoadConfig())
	rbac.Register(public, protected, rbac.Deps{
		DB:      db,
		Binder:  binder,
		Config:  authCfg,
		RDB:     rdb,
		Captcha: captchaSvc,
	})

	settingSvc := app.NewSettingService(db)

	setting.RegisterAdmin(protected, setting.Deps{DB: db, Binder: binder})
	article.Register(protected, article.Deps{DB: db, Binder: binder})
	upload.RegisterAdmin(protected, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	user.Register(protected, user.Deps{DB: db, Binder: binder})
	region.Register(protected, region.Deps{DB: db, Binder: binder})

	health.Register(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.LoadHTMLGlob("./views/*.html")
	r.Static("/assets", "./views/assets")
	r.Static("/uploads", "./uploads")

	r.GET("/favicon.ico", func(c *gin.Context) { c.File("./views/favicon.ico") })
	r.GET("/logo.svg", func(c *gin.Context) { c.File("./views/logo.svg") })
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.NoRoute(func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })

	return r
}

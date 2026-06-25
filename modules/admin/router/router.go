package router

import (
	"net/http"
	"zero-backend/modules/admin/controller"
	adminMiddleware "zero-backend/modules/admin/middleware"
	"zero-backend/modules/article"
	"zero-backend/modules/captcha"
	"zero-backend/modules/health"
	"zero-backend/modules/region"
	"zero-backend/modules/setting"
	"zero-backend/modules/upload"
	"zero-backend/modules/user"

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
	ctrl *controller.Controllers,
	adminMiddlewares *adminMiddleware.Middlewares,
	corsConfig basecfg.CorsConfig,
	authConfig config.AdminAuthConfig,
	db *gorm.DB,
	binder *bind.Binder,
	settingSvc *setting.Service,
	captchaSvc *captcha.Service,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	apiGroup.POST("/login", ctrl.AuthController.Login)
	apiGroup.POST("/refresh-token", ctrl.AuthController.RefreshToken)

	captcha.RegisterWith(apiGroup, binder, captchaSvc)

	apiGroup.Use(middleware.JWTGuard(authConfig.HmacSecret))
	apiGroup.Use(adminMiddlewares.Auth.LoadUser())
	apiGroup.Use(adminMiddlewares.Auth.CheckAPIPermission())

	apiGroup.POST("/logout", ctrl.AuthController.Logout)
	apiGroup.POST("/change-password", ctrl.AuthController.ChangePassword)
	apiGroup.POST("/permissions", ctrl.AuthController.GetPermissions)

	apiGroup.POST("/rbac/menu/list", ctrl.RbacMenuController.List)
	apiGroup.POST("/rbac/menu/create", ctrl.RbacMenuController.Create)
	apiGroup.POST("/rbac/menu/update", ctrl.RbacMenuController.Update)
	apiGroup.POST("/rbac/menu/delete", ctrl.RbacMenuController.Delete)
	apiGroup.POST("/rbac/menu/sync", ctrl.RbacMenuController.Sync)
	apiGroup.POST("/rbac/menu/api/list", ctrl.RbacMenuController.ApiList)
	apiGroup.POST("/rbac/menu/api/save", ctrl.RbacMenuController.ApiSave)
	apiGroup.POST("/rbac/api/list", ctrl.RbacApiController.List)
	apiGroup.POST("/rbac/api/create", ctrl.RbacApiController.Create)
	apiGroup.POST("/rbac/api/update", ctrl.RbacApiController.Update)
	apiGroup.POST("/rbac/api/delete", ctrl.RbacApiController.Delete)
	apiGroup.POST("/rbac/store/list", ctrl.RbacStoreController.List)
	apiGroup.POST("/rbac/store/create", ctrl.RbacStoreController.Create)
	apiGroup.POST("/rbac/store/update", ctrl.RbacStoreController.Update)
	apiGroup.POST("/rbac/store/delete", ctrl.RbacStoreController.Delete)
	apiGroup.POST("/rbac/store/recycle", ctrl.RbacStoreController.Recycle)
	apiGroup.POST("/rbac/store/restore", ctrl.RbacStoreController.Restore)
	apiGroup.POST("/rbac/role/list", ctrl.RbacRoleController.List)
	apiGroup.POST("/rbac/role/create", ctrl.RbacRoleController.Create)
	apiGroup.POST("/rbac/role/update", ctrl.RbacRoleController.Update)
	apiGroup.POST("/rbac/role/delete", ctrl.RbacRoleController.Delete)
	apiGroup.POST("/rbac/role/set-menus", ctrl.RbacRoleController.SetMenus)
	apiGroup.POST("/rbac/user/list", ctrl.RbacUserController.List)
	apiGroup.POST("/rbac/user/create", ctrl.RbacUserController.Create)
	apiGroup.POST("/rbac/user/update", ctrl.RbacUserController.Update)
	apiGroup.POST("/rbac/user/delete", ctrl.RbacUserController.Delete)
	apiGroup.POST("/rbac/user/set-roles", ctrl.RbacUserController.SetRoles)
	apiGroup.POST("/rbac/user/reset-password", ctrl.RbacUserController.ResetPassword)

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

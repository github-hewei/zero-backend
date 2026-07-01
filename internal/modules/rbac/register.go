package rbac

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB      *gorm.DB
	Binder  *bind.Binder
	Config  Config
	RDB     *redis.Client
	Captcha CaptchaVerifier
}

// PlatformDeps 平台端模块依赖（不含 Captcha，不挂 RBAC 中间件）
type PlatformDeps struct {
	DB     *gorm.DB
	Binder *bind.Binder
	Config Config
	RDB    *redis.Client
}

func buildAll(deps Deps) (*handler, *AuthMiddleware) {
	userRepo := NewRbacUserRepository(deps.DB)
	userRoleRepo := NewRbacUserRoleRepository(deps.DB)
	menuRepo := NewRbacMenuRepository(deps.DB)
	menuApiRepo := NewRbacMenuApiRepository(deps.DB)
	apiRepo := NewRbacApiRepository(deps.DB)
	roleRepo := NewRbacRoleRepository(deps.DB)
	roleMenuRepo := NewRbacRoleMenuRepository(deps.DB)
	storeRepo := NewRbacStoreRepository(deps.DB)

	authServ := NewAuthService(
		userRepo, apiRepo, roleRepo, menuRepo,
		roleMenuRepo, userRoleRepo, menuApiRepo,
		deps.Config, deps.RDB, deps.Captcha,
	)

	authMid := NewAuthMiddleware(deps.Config, authServ)

	menuServ := NewRbacMenuService(menuRepo, menuApiRepo, deps.DB)
	apiServ := NewRbacApiService(apiRepo)
	roleServ := NewRbacRoleService(roleRepo, roleMenuRepo, deps.DB)
	userServ := NewRbacUserService(deps.DB, userRepo, userRoleRepo)
	storeServ := NewRbacStoreService(storeRepo)

	h := newHandler(deps.Binder, authServ, deps.Config, menuServ, apiServ, roleServ, userServ, storeServ, authMid)

	return h, authMid
}

// Register 注册 rbac 模块路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT 中间件。
// admin 端仅管理自身租户的角色和用户，store/menu/api 管理已迁移到平台端。
func Register(public, protected *gin.RouterGroup, deps Deps) {
	h, authMid := buildAll(deps)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(deps.Config.HmacSecret))
	protected.Use(authMid.LoadUser())
	protected.Use(authMid.CheckAPIPermission())

	protected.POST("/logout", h.logout)
	protected.POST("/change-password", h.changePassword)
	protected.POST("/permissions", h.permissions)

	protected.POST("/rbac/role/list", h.roleList)
	protected.POST("/rbac/role/create", h.roleCreate)
	protected.POST("/rbac/role/update", h.roleUpdate)
	protected.POST("/rbac/role/delete", h.roleDelete)
	protected.POST("/rbac/role/set-menus", h.roleSetMenus)

	protected.POST("/rbac/user/list", h.userList)
	protected.POST("/rbac/user/create", h.userCreate)
	protected.POST("/rbac/user/update", h.userUpdate)
	protected.POST("/rbac/user/delete", h.userDelete)
	protected.POST("/rbac/user/set-roles", h.userSetRoles)
	protected.POST("/rbac/user/reset-password", h.userResetPassword)
}

// RegisterPlatform 注册平台端路由。不挂载 RBAC 权限中间件，由平台端的 RequireRole 中间件控制访问权限。
func RegisterPlatform(rg *gin.RouterGroup, deps PlatformDeps) {
	fullDeps := Deps{
		DB:     deps.DB,
		Binder: deps.Binder,
		Config: deps.Config,
		RDB:    deps.RDB,
	}
	h, authMid := buildAll(fullDeps)

	rg.Use(middleware.JWTGuard(deps.Config.HmacSecret))
	rg.Use(authMid.LoadUser())

	rg.POST("/rbac/menu/list", h.menuList)
	rg.POST("/rbac/menu/create", h.menuCreate)
	rg.POST("/rbac/menu/update", h.menuUpdate)
	rg.POST("/rbac/menu/delete", h.menuDelete)
	rg.POST("/rbac/menu/sync", h.menuSync)
	rg.POST("/rbac/menu/api/list", h.menuApiList)
	rg.POST("/rbac/menu/api/save", h.menuApiSave)

	rg.POST("/rbac/api/list", h.apiList)
	rg.POST("/rbac/api/create", h.apiCreate)
	rg.POST("/rbac/api/update", h.apiUpdate)
	rg.POST("/rbac/api/delete", h.apiDelete)

	rg.POST("/rbac/store/list", h.storeList)
	rg.POST("/rbac/store/create", h.storeCreate)
	rg.POST("/rbac/store/update", h.storeUpdate)
	rg.POST("/rbac/store/delete", h.storeDelete)
	rg.POST("/rbac/store/recycle", h.storeRecycle)
	rg.POST("/rbac/store/restore", h.storeRestore)

	rg.POST("/rbac/role/list", h.roleList)
	rg.POST("/rbac/role/create", h.roleCreate)
	rg.POST("/rbac/role/update", h.roleUpdate)
	rg.POST("/rbac/role/delete", h.roleDelete)
	rg.POST("/rbac/role/set-menus", h.roleSetMenus)

	rg.POST("/rbac/user/list", h.userList)
	rg.POST("/rbac/user/create", h.userCreate)
	rg.POST("/rbac/user/update", h.userUpdate)
	rg.POST("/rbac/user/delete", h.userDelete)
	rg.POST("/rbac/user/reset-password", h.userResetPassword)
}

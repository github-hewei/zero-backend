package controller

// Controllers 控制器集合
type Controllers struct {
	AuthController            *AuthController
	CaptchaController         *CaptchaController
	RbacMenuController        *RbacMenuController
	RbacApiController         *RbacApiController
	RbacRoleController        *RbacRoleController
	RbacUserController        *RbacUserController
	RbacStoreController       *RbacStoreController
	UserController            *UserController
	SettingController         *SettingController
	SettingDefaultController  *SettingDefaultController
	RegionController          *RegionController
	HealthController          *HealthController
}

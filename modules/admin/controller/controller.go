package controller

// Controllers 控制器集合
type Controllers struct {
	AuthController            *AuthController
	RbacMenuController        *RbacMenuController
	RbacApiController         *RbacApiController
	RbacRoleController        *RbacRoleController
	RbacUserController        *RbacUserController
	RbacStoreController       *RbacStoreController
	UploadGroupController     *UploadGroupController
	UploadFileController      *UploadFileController
	UserController            *UserController
	SettingController         *SettingController
	SettingDefaultController  *SettingDefaultController
	ArticleCategoryController *ArticleCategoryController
	ArticleController         *ArticleController
	RegionController          *RegionController
	HealthController          *HealthController
}

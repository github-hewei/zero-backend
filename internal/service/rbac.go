package service

import (
	"context"
	"zero-backend/internal/apperror"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
	"zero-backend/pkg/helper"

	"gorm.io/gorm"
)

// RbacUserService 用户服务
type RbacUserService struct {
	db           *gorm.DB
	repo         *repository.RbacUserRepository
	userRoleRepo *repository.RbacUserRoleRepository
}

// NewRbacUserService 创建用户服务
func NewRbacUserService(db *gorm.DB, repo *repository.RbacUserRepository, userRoleRepo *repository.RbacUserRoleRepository) *RbacUserService {
	return &RbacUserService{
		db:           db,
		repo:         repo,
		userRoleRepo: userRoleRepo,
	}
}

// FindList 获取用户列表
func (s *RbacUserService) FindList(ctx context.Context, req *dto.RbacUserListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.RbacUser{},
		Total: 0,
	}

	filter := &repository.RbacUserFilterField{
		Username: req.Username,
		RealName: req.RealName,
		StoreId:  req.StoreId,
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "sort", Sort: "asc"},
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询用户数量失败")
	}

	if total == 0 {
		return result, nil
	}

	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders,
		repository.WithScopes(nil), // 覆盖掉默认的数据隔离条件，实现用默认的数据隔离
		repository.WithPreloads("RbacUserRole.RbacRole"),
	)

	if err != nil {
		return nil, apperror.NewSystemError(err, "查询用户列表失败")
	}

	result.List = list

	return result, nil
}

// Create 创建用户
func (s *RbacUserService) Create(ctx context.Context, req *dto.RbacUserCreateRequest) error {
	if err := s.checkUsername(ctx, req.Username, req.StoreId); err != nil {
		return err
	}

	// 密码加密
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return apperror.NewSystemError(err, "密码加密失败")
	}

	item := &model.RbacUser{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		IsSuper:  req.IsSuper,
		StoreId:  req.StoreId,
		Sort:     req.Sort,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建用户失败")
	}

	return nil
}

// Update 更新用户
func (s *RbacUserService) Update(ctx context.Context, req *dto.RbacUserUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询用户失败")
	}

	if item.ID == 0 || item.StoreId != req.StoreId {
		return apperror.NewUserError("用户不存在")
	}

	if item.Username != req.Username {
		if err := s.checkUsername(ctx, req.Username, req.StoreId); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"username":  req.Username,
		"real_name": req.RealName,
		"store_id":  req.StoreId,
		"is_super":  req.IsSuper,
		"sort":      req.Sort,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新用户失败")
	}

	return nil
}

// Delete 删除用户
func (s *RbacUserService) Delete(ctx context.Context, req *dto.RbacUserDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询用户失败")
	}

	if item.ID == 0 || item.StoreId != req.StoreId {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除用户失败")
	}

	return nil
}

// SetRoles 设置用户角色
func (s *RbacUserService) SetRoles(ctx context.Context, req *dto.RbacUserRoleSetRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 验证用户存在
		user, err := s.repo.FindOne(ctx, req.UserID, repository.WithTx[*repository.QueryConfig](tx), repository.WithScopes(nil))
		if err != nil {
			return apperror.NewSystemError(err, "查询用户失败")
		}

		if user.ID == 0 || user.StoreId != req.StoreId {
			return apperror.NewUserError("用户不存在")
		}

		// 获取现有角色
		filter := &repository.RbacUserRoleFilterField{UserId: user.ID}
		existingRoles, err := s.userRoleRepo.FindAll(ctx, filter, nil, nil,
			repository.WithScopes(nil),
			repository.WithTx[*repository.QueryConfig](tx),
		)
		if err != nil {
			return apperror.NewSystemError(err, "查询用户角色失败")
		}

		// 计算差异
		existingMap := make(map[uint32]bool)
		for _, role := range existingRoles {
			existingMap[role.RoleId] = true
		}

		newMap := make(map[uint32]bool)
		for _, id := range req.RoleIDs {
			newMap[id] = true
		}

		// 删除不再需要的角色
		deleteIds := make([]uint32, 0)
		for _, role := range existingRoles {
			if !newMap[role.RoleId] {
				deleteIds = append(deleteIds, role.ID)
			}
		}
		if len(deleteIds) > 0 {
			if err := s.userRoleRepo.Delete(ctx, deleteIds, repository.WithTx[*repository.DeleteConfig](tx)); err != nil {
				return apperror.NewSystemError(err, "删除用户角色失败")
			}
		}

		// 添加新角色
		createUserRoles := make([]*model.RbacUserRole, 0)
		for _, roleID := range req.RoleIDs {
			if !existingMap[roleID] {
				createUserRoles = append(createUserRoles, &model.RbacUserRole{
					UserId:  user.ID,
					RoleId:  roleID,
					StoreId: req.StoreId,
				})
			}
		}

		if err := s.userRoleRepo.CreateBatch(ctx, createUserRoles, repository.WithTx[*repository.CreateConfig](tx)); err != nil {
			return apperror.NewSystemError(err, "创建用户角色失败")
		}

		return nil
	})
}

// checkUsername 检查用户名是否已存在
func (s *RbacUserService) checkUsername(ctx context.Context, username string, storeId uint32) error {
	filter := &repository.RbacUserUsernameFilterField{Username: username, StoreId: storeId}
	item, err := s.repo.FindOne(ctx, filter)

	if err != nil {
		return apperror.NewSystemError(err, "查询用户名失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("用户名已存在")
	}

	return nil
}

// ResetPassword 重置用户密码
func (s *RbacUserService) ResetPassword(ctx context.Context, req *dto.RbacUserResetPasswordRequest) (string, error) {
	user, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return "", apperror.NewSystemError(err, "查询用户失败")
	}

	if user.ID == 0 || user.StoreId != req.StoreId {
		return "", apperror.NewUserError("用户不存在")
	}

	// 生成12位包含特殊字符的随机密码
	newPassword := helper.RandomStringWithSymbols(12)
	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return "", apperror.NewSystemError(err, "密码加密失败")
	}

	// 更新密码
	updateData := map[string]any{
		"password": hashedPassword,
	}

	if err := s.repo.Updates(ctx, user, updateData); err != nil {
		return "", apperror.NewSystemError(err, "重置密码失败")
	}

	return newPassword, nil
}

// RbacMenuService 菜单服务
type RbacMenuService struct {
	repo    *repository.RbacMenuRepository
	apiRepo *repository.RbacMenuApiRepository
	Db      *gorm.DB
}

// NewRbacMenuService 创建菜单服务
func NewRbacMenuService(repo *repository.RbacMenuRepository, apiRepo *repository.RbacMenuApiRepository, db *gorm.DB) *RbacMenuService {
	return &RbacMenuService{
		repo:    repo,
		apiRepo: apiRepo,
		Db:      db,
	}
}

// FindTreeList 获取菜单树
func (s *RbacMenuService) FindTreeList(ctx context.Context) ([]*model.RbacMenu, error) {
	list, err := s.repo.FindAll(ctx, nil, nil, nil, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询菜单列表失败")
	}

	if len(list) > 0 {
		rbacMenuList := model.RbacMenuList{}
		for _, menu := range list {
			rbacMenuList = append(rbacMenuList, menu)
		}

		list = rbacMenuList.Tree()
	}

	return list, nil
}

// Create 创建菜单
func (s *RbacMenuService) Create(ctx context.Context, req *dto.RbacMenuCreateRequest) error {
	if err := s.checkName(ctx, req.Name); err != nil {
		return err
	}

	item := &model.RbacMenu{
		Type:       req.Type,
		Name:       req.Name,
		ParentId:   req.ParentId,
		Path:       req.Path,
		Sort:       req.Sort,
		IsPage:     req.IsPage,
		ModuleKey:  req.ModuleKey,
		ActionMark: req.ActionMark,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建菜单失败")
	}

	return nil
}

// Update 更新菜单
func (s *RbacMenuService) Update(ctx context.Context, req *dto.RbacMenuUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询菜单失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if item.Name != req.Name {
		if err := s.checkName(ctx, req.Name); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"type":        req.Type,
		"name":        req.Name,
		"path":        req.Path,
		"sort":        req.Sort,
		"module_key":  req.ModuleKey,
		"parent_id":   req.ParentId,
		"is_page":     req.IsPage,
		"action_mark": req.ActionMark,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新菜单失败")
	}

	return nil
}

// Delete 删除菜单
func (s *RbacMenuService) Delete(ctx context.Context, req *dto.RbacMenuDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询菜单失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除菜单失败")
	}

	return nil
}

// Sync 同步菜单
func (s *RbacMenuService) Sync(ctx context.Context, req []dto.RbacMenuSyncRequest) error {
	// 开启事务
	return s.Db.Transaction(func(tx *gorm.DB) error {
		list, err := s.repo.FindAll(ctx, &repository.RbacMenuFilterField{Type: 10}, nil, nil,
			repository.WithTx[*repository.QueryConfig](tx),
			repository.WithScopes(nil),
		)
		if err != nil {
			return apperror.NewSystemError(err, "查询菜单列表失败")
		}

		menuMap := map[string]*model.RbacMenu{}
		for _, item := range list {
			menuMap[item.Path] = item
		}

		if err := s.SyncMenuList(ctx, req, 0, menuMap, tx); err != nil {
			return err
		}

		for _, item := range menuMap {
			if err := s.repo.Delete(ctx, item.ID, repository.WithTx[*repository.DeleteConfig](tx)); err != nil {
				return apperror.NewSystemError(err, "删除菜单失败")
			}
		}
		return nil
	})
}

// SyncMenuList 递归同步菜单
func (s *RbacMenuService) SyncMenuList(ctx context.Context, req []dto.RbacMenuSyncRequest, parentId uint32, menuMap map[string]*model.RbacMenu, tx *gorm.DB) error {
	for _, item := range req {
		menu, exists := menuMap[item.Path]
		if exists {
			updateData := map[string]any{
				"module_key": item.ModuleKey,
				"sort":       item.Sort,
				"type":       item.Type,
				"name":       item.Name,
				"is_page":    item.IsPage,
				"parent_id":  parentId,
			}

			if err := s.repo.Updates(ctx, menu, updateData, repository.WithTx[*repository.UpdateConfig](tx)); err != nil {
				return apperror.NewSystemError(err, "更新菜单失败")
			}
			delete(menuMap, item.Path)

		} else {
			menu = &model.RbacMenu{
				Type:      item.Type,
				Name:      item.Name,
				Path:      item.Path,
				IsPage:    item.IsPage,
				ModuleKey: item.ModuleKey,
				ParentId:  parentId,
				Sort:      item.Sort,
			}

			if err := s.repo.Create(ctx, menu, repository.WithTx[*repository.CreateConfig](tx)); err != nil {
				return apperror.NewSystemError(err, "创建菜单失败")
			}
		}

		if len(item.Children) > 0 {
			if err := s.SyncMenuList(ctx, item.Children, menu.ID, menuMap, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkName 检查菜单名称
func (s *RbacMenuService) checkName(ctx context.Context, name string) error {
	filter := &repository.RbacMenuFilterField{Name: name}
	item, err := s.repo.FindOne(ctx, filter, repository.WithScopes(nil))

	if err != nil {
		return apperror.NewSystemError(err, "查询菜单名失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("菜单名已存在")
	}

	return nil
}

// FindApiList 查询菜单关联的api
func (s *RbacMenuService) FindApiList(ctx context.Context, req *dto.RbacMenuApiListRequest) (*dto.MenuApiRelationResponse, error) {
	filter := &repository.RbacMenuApiFilterField{MenuId: req.MenuID}
	list, err := s.apiRepo.FindAll(ctx, filter, nil, nil, repository.WithScopes(nil))

	if err != nil {
		return nil, apperror.NewSystemError(err, "查询菜单权限失败")
	}

	apiIds := make([]uint32, 0)
	for _, item := range list {
		apiIds = append(apiIds, item.ApiId)
	}

	return &dto.MenuApiRelationResponse{
		ApiIds: apiIds,
		MenuId: req.MenuID,
	}, nil
}

// SaveApiList 保存菜单权限
func (s *RbacMenuService) SaveApiList(ctx context.Context, req *dto.RbacMenuApiSaveRequest) error {
	// 开启事务
	return s.Db.Transaction(func(tx *gorm.DB) error {
		filter := &repository.RbacMenuApiFilterField{MenuId: req.MenuID}
		list, err := s.apiRepo.FindAll(ctx, filter, nil, nil,
			repository.WithTx[*repository.QueryConfig](tx),
			repository.WithScopes(nil),
		)

		if err != nil {
			return apperror.NewSystemError(err, "查询菜单权限失败")
		}

		apiIds := make(map[uint32]bool)
		for _, item := range list {
			apiIds[item.ApiId] = true
		}

		for _, apiId := range req.ApiIds {
			if !apiIds[apiId] {
				item := &model.RbacMenuApi{
					MenuId: req.MenuID,
					ApiId:  apiId,
				}
				err := s.apiRepo.Create(ctx, item, repository.WithTx[*repository.CreateConfig](tx))
				if err != nil {
					return apperror.NewSystemError(err, "保存菜单权限失败")
				}
			} else {
				delete(apiIds, apiId)
			}
		}

		for apiId := range apiIds {
			filter := &repository.RbacMenuApiFilterField{ApiId: apiId, MenuId: req.MenuID}
			err := s.apiRepo.Delete(ctx, filter, repository.WithTx[*repository.DeleteConfig](tx))

			if err != nil {
				return apperror.NewSystemError(err, "删除菜单权限失败")
			}
		}

		return nil
	})
}

// RbacRoleService 角色服务
type RbacRoleService struct {
	db           *gorm.DB
	repo         *repository.RbacRoleRepository
	roleMenuRepo *repository.RbacRoleMenuRepository
}

// NewRbacRoleService 创建角色服务
func NewRbacRoleService(repo *repository.RbacRoleRepository, roleMenuRepo *repository.RbacRoleMenuRepository, db *gorm.DB) *RbacRoleService {
	return &RbacRoleService{
		db:           db,
		repo:         repo,
		roleMenuRepo: roleMenuRepo,
	}
}

// FindTreeList 获取角色树形列表
func (s *RbacRoleService) FindTreeList(ctx context.Context, req *dto.RbacRoleListRequest) ([]*model.RbacRole, error) {
	filter := &repository.RbacRoleFilterField{
		StoreId:  req.StoreId,
		RoleName: req.RoleName,
	}

	list, err := s.repo.FindAll(ctx, filter, nil, nil,
		repository.WithScopes(nil), // 覆盖掉默认的数据隔离条件，实现用默认的数据隔离
		repository.WithPreloads("RbacRoleMenu.RbacMenu"),
	)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询角色列表失败")
	}

	if len(list) > 0 {
		rbacRoleList := model.RbacRoleList{}
		for _, role := range list {
			rbacRoleList = append(rbacRoleList, role)
		}
		list = rbacRoleList.Tree()
	}

	return list, nil
}

// Create 创建角色
func (s *RbacRoleService) Create(ctx context.Context, req *dto.RbacRoleCreateRequest) error {
	if err := s.checkName(ctx, req.RoleName, req.StoreId); err != nil {
		return err
	}

	if req.ParentId > 0 {
		if err := s.checkParent(ctx, req.ParentId, req.StoreId); err != nil {
			return err
		}
	}

	item := &model.RbacRole{
		RoleName: req.RoleName,
		ParentId: req.ParentId,
		Sort:     req.Sort,
		StoreId:  req.StoreId,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建角色失败")
	}

	return nil
}

// Update 更新角色
func (s *RbacRoleService) Update(ctx context.Context, req *dto.RbacRoleUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询角色失败")
	}

	if item.ID == 0 || item.StoreId != req.StoreId {
		return apperror.NewUserError("找不到此记录")
	}

	if item.RoleName != req.RoleName {
		if err := s.checkName(ctx, req.RoleName, req.StoreId); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"role_name": req.RoleName,
		"parent_id": req.ParentId,
		"sort":      req.Sort,
		"store_id":  req.StoreId,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新角色失败")
	}

	return nil
}

// Delete 删除角色
func (s *RbacRoleService) Delete(ctx context.Context, req *dto.RbacRoleDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询角色失败")
	}

	if item.ID == 0 || item.StoreId != req.StoreId {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除角色失败")
	}

	return nil
}

// SetMenus 设置角色菜单
func (s *RbacRoleService) SetMenus(ctx context.Context, req *dto.RbacRoleMenuSetRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 验证角色存在
		role, err := s.repo.FindOne(ctx, req.RoleID, repository.WithTx[*repository.QueryConfig](tx), repository.WithScopes(nil))

		if err != nil {
			return apperror.NewSystemError(err, "查询角色失败")
		}

		if role.ID == 0 || role.StoreId != req.StoreId {
			return apperror.NewUserError("角色不存在")
		}

		filter := &repository.RbacRoleMenuFilterField{RoleId: role.ID}
		existingMenus, err := s.roleMenuRepo.FindAll(ctx, filter, nil, nil,
			repository.WithTx[*repository.QueryConfig](tx),
			repository.WithScopes(nil),
		)

		if err != nil {
			return apperror.NewSystemError(err, "查询角色菜单失败")
		}

		// 计算差异
		existingMap := make(map[uint32]bool)
		for _, menu := range existingMenus {
			existingMap[menu.MenuId] = true
		}

		newMap := make(map[uint32]bool)
		for _, id := range req.MenuIDs {
			newMap[id] = true
		}

		// 删除不再需要的菜单
		deleteIds := make([]uint32, 0)
		for _, menu := range existingMenus {
			if !newMap[menu.MenuId] {
				deleteIds = append(deleteIds, menu.ID)
			}
		}
		if len(deleteIds) > 0 {
			if err := s.roleMenuRepo.Delete(ctx, deleteIds, repository.WithTx[*repository.DeleteConfig](tx)); err != nil {
				return err
			}
		}

		// 添加新菜单
		roleMenus := make([]*model.RbacRoleMenu, 0)
		for _, menuID := range req.MenuIDs {
			if !existingMap[menuID] {
				roleMenus = append(roleMenus, &model.RbacRoleMenu{
					RoleId:  role.ID,
					MenuId:  menuID,
					StoreId: req.StoreId,
				})
			}
		}

		if err := s.roleMenuRepo.CreateBatch(ctx, roleMenus, repository.WithTx[*repository.CreateConfig](tx)); err != nil {
			return err
		}

		return nil
	})
}

// checkName 检查角色名称
func (s *RbacRoleService) checkName(ctx context.Context, name string, StoreId uint32) error {
	item, err := s.repo.FindByName(ctx, name, StoreId)
	if err != nil {
		return apperror.NewSystemError(err, "查询角色名失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("角色名已存在")
	}

	return nil
}

// checkParent 检查父级角色
func (s *RbacRoleService) checkParent(ctx context.Context, parentId uint32, StoreId uint32) error {
	parent, err := s.repo.FindOne(ctx, parentId)
	if err != nil {
		return apperror.NewSystemError(err, "查询父级角色失败")
	}

	if parent.ID == 0 {
		return apperror.NewUserError("父级角色不存在")
	}

	if parent.StoreId != StoreId {
		return apperror.NewUserError("父级角色不属于当前企业")
	}

	return nil
}

// RbacApiService 接口服务
type RbacApiService struct {
	repo *repository.RbacApiRepository
}

// NewRbacApiService 创建接口服务
func NewRbacApiService(repo *repository.RbacApiRepository) *RbacApiService {
	return &RbacApiService{repo: repo}
}

// FindTreeList 获取接口树列表
func (s *RbacApiService) FindTreeList(ctx context.Context) ([]*model.RbacApi, error) {
	list, err := s.repo.FindAll(ctx, nil, nil, nil, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询接口列表失败")
	}

	if len(list) > 0 {
		rbacApiList := model.RbacApiList{}
		for _, api := range list {
			rbacApiList = append(rbacApiList, api)
		}

		list = rbacApiList.Tree()
	}

	return list, nil
}

// Create 创建接口
func (s *RbacApiService) Create(ctx context.Context, req *dto.RbacApiCreateRequest) error {
	if err := s.checkName(ctx, req.Name); err != nil {
		return err
	}

	item := &model.RbacApi{
		Name:     req.Name,
		ParentId: req.ParentId,
		Url:      req.Url,
		Sort:     req.Sort,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建接口失败")
	}

	return nil
}

// Update 更新接口
func (s *RbacApiService) Update(ctx context.Context, req *dto.RbacApiUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询接口失败")
	}

	if item.Name != req.Name {
		if err := s.checkName(ctx, req.Name); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"name":      req.Name,
		"parent_id": req.ParentId,
		"url":       req.Url,
		"sort":      req.Sort,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新接口失败")
	}

	return nil
}

// Delete 删除接口
func (s *RbacApiService) Delete(ctx context.Context, req *dto.RbacApiDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询接口失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除接口失败")
	}

	return nil
}

// checkName 检查接口名称
func (s *RbacApiService) checkName(ctx context.Context, name string) error {
	filter := &repository.RbacApiFilterField{Name: name}
	item, err := s.repo.FindOne(ctx, filter, repository.WithScopes(nil))

	if err != nil {
		return apperror.NewSystemError(err, "查询接口名失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("接口名已存在")
	}

	return nil
}

// RbacStoreService 企业服务
type RbacStoreService struct {
	repo *repository.RbacStoreRepository
}

// NewRbacStoreService 创建企业服务
func NewRbacStoreService(repo *repository.RbacStoreRepository) *RbacStoreService {
	return &RbacStoreService{repo: repo}
}

// FindList 获取企业列表
func (s *RbacStoreService) FindList(ctx context.Context, req *dto.RbacStoreListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.RbacStore{},
		Total: 0,
	}

	filter := &repository.RbacStoreFilterField{
		Name:      req.Name,
		IsRecycle: req.IsRecycle,
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "sort", Sort: "asc"},
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询企业数量失败")
	}

	if total == 0 {
		return result, nil
	}

	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询企业列表失败")
	}

	result.List = list

	return result, nil
}

// Create 创建企业
func (s *RbacStoreService) Create(ctx context.Context, req *dto.RbacStoreCreateRequest) error {
	if err := s.checkName(ctx, req.Name); err != nil {
		return err
	}

	item := &model.RbacStore{
		Name:         req.Name,
		ShortName:    req.ShortName,
		Contact:      req.Contact,
		ContactPhone: req.ContactPhone,
		Description:  req.Description,
		LogoImageID:  req.LogoImageID,
		Sort:         req.Sort,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建企业失败")
	}

	return nil
}

// checkName 检查企业名称
func (s *RbacStoreService) checkName(ctx context.Context, name string) error {
	filter := &repository.RbacStoreNameFilterField{Name: name}
	item, err := s.repo.FindOne(ctx, filter, repository.WithScopes(nil))

	if err != nil {
		return apperror.NewSystemError(err, "查询企业名失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("企业名已存在")
	}

	return nil
}

// Update 更新企业信息
func (s *RbacStoreService) Update(ctx context.Context, req *dto.RbacStoreUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询企业失败")
	}

	if item.Name != req.Name {
		if err := s.checkName(ctx, req.Name); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"contact":       req.Contact,
		"contact_phone": req.ContactPhone,
		"description":   req.Description,
		"name":          req.Name,
		"short_name":    req.ShortName,
		"sort":          req.Sort,
		"logo_image_id": req.LogoImageID,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新企业失败")
	}

	return nil
}

// Delete 删除企业信息
func (s *RbacStoreService) Delete(ctx context.Context, req *dto.RbacStoreDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询企业失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除企业失败")
	}

	return nil
}

// Recycle 将企业移入回收站
func (s *RbacStoreService) Recycle(ctx context.Context, req *dto.RbacStoreDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID)
	if err != nil {
		return apperror.NewSystemError(err, "查询企业失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	updateData := map[string]any{
		"is_recycle": 1,
	}
	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "移入回收站失败")
	}

	return nil
}

// Restore 将企业从回收站恢复
func (s *RbacStoreService) Restore(ctx context.Context, req *dto.RbacStoreDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID)
	if err != nil {
		return apperror.NewSystemError(err, "查询企业失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	updateData := map[string]any{
		"is_recycle": 0,
	}
	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "从回收站恢复失败")
	}

	return nil
}

package repository

import (
	"context"
	"errors"
	"zero-backend/internal/model"
	"zero-backend/pkg/helper"

	"gorm.io/gorm"
)

// RbacUserFilterField 用户表过滤字段
type RbacUserFilterField struct {
	Username string
	RealName string
	StoreId  uint32
}

// Apply 应用筛选条件
func (f *RbacUserFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Username != "" {
		db = db.Where("username like ?", helper.SafeLikeString(f.Username)+"%")
	}

	if f.RealName != "" {
		db = db.Where("real_name like ?", helper.SafeLikeString(f.RealName)+"%")
	}

	if f.StoreId > 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	return db
}

// RbacUserUsernameFilterField 用户名筛选条件
type RbacUserUsernameFilterField struct {
	Username string
}

// Apply 筛选条件
func (f *RbacUserUsernameFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Username != "" {
		db = db.Where("username = ?", f.Username)
	}

	return db
}

// RbacUserRepository 用户数据操作
type RbacUserRepository struct {
	*BaseRepository[model.RbacUser]
}

// NewRbacUserRepository 创建用户数据操作
func NewRbacUserRepository(db *gorm.DB) *RbacUserRepository {
	baseRepo := NewBaseRepository[model.RbacUser](db)
	return &RbacUserRepository{
		BaseRepository: baseRepo,
	}
}

// RbacUserRoleFilterField 查询字段
type RbacUserRoleFilterField struct {
	UserId uint32
}

// Apply 应用过滤条件
func (f *RbacUserRoleFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.UserId > 0 {
		db = db.Where("user_id = ?", f.UserId)
	}

	return db
}

// RbacUserRoleRepository 用户关联角色数据操作
type RbacUserRoleRepository struct {
	*BaseRepository[model.RbacUserRole]
}

// NewRbacUserRoleRepository 创建用户关联角色数据操作
func NewRbacUserRoleRepository(db *gorm.DB) *RbacUserRoleRepository {
	baseRepo := NewBaseRepository[model.RbacUserRole](db)
	return &RbacUserRoleRepository{
		BaseRepository: baseRepo,
	}
}

// RbacMenuFilterField 菜单表表过滤字段
type RbacMenuFilterField struct {
	IDs  []uint32
	Type int8
	Name string
}

// Apply 应用筛选条件
func (f *RbacMenuFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if len(f.IDs) > 0 {
		db = db.Where("id IN ?", f.IDs)
	}

	if f.Type != 0 {
		db = db.Where("type = ?", f.Type)
	}

	if f.Name != "" {
		db = db.Where("name = ?", f.Name)
	}

	return db
}

// RbacMenuRepository 菜单数据操作
type RbacMenuRepository struct {
	*BaseRepository[model.RbacMenu]
}

// NewRbacMenuRepository 创建菜单数据操作
func NewRbacMenuRepository(db *gorm.DB) *RbacMenuRepository {
	baseRepo := NewBaseRepository[model.RbacMenu](db)
	return &RbacMenuRepository{
		BaseRepository: baseRepo,
	}
}

// RbacMenuApiFilterField 菜单关联接口表过滤字段
type RbacMenuApiFilterField struct {
	MenuId  uint32
	ApiId   uint32
	MenuIDs []uint32
}

// Apply 应用筛选条件
func (f *RbacMenuApiFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.MenuId > 0 {
		db = db.Where("menu_id = ?", f.MenuId)
	}

	if f.ApiId > 0 {
		db = db.Where("api_id = ?", f.ApiId)
	}

	if len(f.MenuIDs) > 0 {
		db = db.Where("menu_id in (?)", f.MenuIDs)
	}

	return db
}

// RbacMenuApiRepository 菜单关联接口数据操作
type RbacMenuApiRepository struct {
	*BaseRepository[model.RbacMenuApi]
}

// NewRbacMenuApiRepository 创建菜单关联接口数据操作
func NewRbacMenuApiRepository(db *gorm.DB) *RbacMenuApiRepository {
	baseRepo := NewBaseRepository[model.RbacMenuApi](db)
	return &RbacMenuApiRepository{
		BaseRepository: baseRepo,
	}
}

// RbacRoleFilterField 角色表过滤字段
type RbacRoleFilterField struct {
	IDs      []uint32
	StoreId  uint32
	RoleName string
}

// Apply 应用筛选条件
func (f *RbacRoleFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if len(f.IDs) > 0 {
		db = db.Where("id in ?", f.IDs)
	}

	if f.RoleName != "" {
		db = db.Where("role_name like ?", helper.SafeLikeString(f.RoleName)+"%")
	}

	if f.StoreId > 0 {
		db = db.Where("store_id = ?", f.StoreId)
	}

	return db
}

// RbacRoleRepository 角色数据操作
type RbacRoleRepository struct {
	*BaseRepository[model.RbacRole]
}

// NewRbacRoleRepository 创建角色数据操作
func NewRbacRoleRepository(db *gorm.DB) *RbacRoleRepository {
	baseRepo := NewBaseRepository[model.RbacRole](db)
	return &RbacRoleRepository{
		BaseRepository: baseRepo,
	}
}

// FindByName 根据名称获取角色
func (r *RbacRoleRepository) FindByName(ctx context.Context, name string, StoreId uint32) (*model.RbacRole, error) {
	item := &model.RbacRole{}
	err := r.Db.WithContext(ctx).Where("role_name = ? AND store_id = ?", name, StoreId).First(&item).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return item, nil
}

// RbacRoleMenuFilterField 角色菜单过滤字段
type RbacRoleMenuFilterField struct {
	RoleId uint32
}

// Apply 应用筛选条件
func (f *RbacRoleMenuFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.RoleId > 0 {
		db = db.Where("role_id = ?", f.RoleId)
	}

	return db
}

// RbacRoleMenuRepository 角色关联菜单数据操作
type RbacRoleMenuRepository struct {
	*BaseRepository[model.RbacRoleMenu]
}

// NewRbacRoleMenuRepository 创建角色关联菜单数据操作
func NewRbacRoleMenuRepository(db *gorm.DB) *RbacRoleMenuRepository {
	baseRepo := NewBaseRepository[model.RbacRoleMenu](db)
	return &RbacRoleMenuRepository{
		BaseRepository: baseRepo,
	}
}

// RbacApiFilterField 筛选字段
type RbacApiFilterField struct {
	Name string
}

func (f *RbacApiFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Name != "" {
		db = db.Where("name = ?", f.Name)
	}

	return db
}

// RbacApiRepository 接口数据操作
type RbacApiRepository struct {
	*BaseRepository[model.RbacApi]
}

// NewRbacApiRepository 创建接口数据操作
func NewRbacApiRepository(db *gorm.DB) *RbacApiRepository {
	baseRepo := NewBaseRepository[model.RbacApi](db)
	return &RbacApiRepository{
		BaseRepository: baseRepo,
	}
}

// GetAPIByPath 根据路径获取API
func (r *RbacApiRepository) GetAPIByPath(ctx context.Context, path string) (*model.RbacApi, error) {
	item := new(model.RbacApi)
	err := r.Db.WithContext(ctx).Where("url = ?", path).First(item).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return item, nil
}

// RbacStoreFilterField 企业表过滤字段
type RbacStoreFilterField struct {
	Name      string
	IsRecycle int8
}

// Apply 应用筛选条件
func (f *RbacStoreFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Name != "" {
		db = db.Where("name like ?", helper.SafeLikeString(f.Name)+"%")
	}

	if f.IsRecycle != -1 {
		db = db.Where("is_recycle = ?", f.IsRecycle)
	}

	return db
}

// RbacStoreNameFilterField 企业表过滤字段
type RbacStoreNameFilterField struct {
	Name string
}

// Apply 应用筛选条件
func (f *RbacStoreNameFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Name != "" {
		db = db.Where("name = ?", f.Name)
	}

	return db
}

// RbacStoreRepository 企业数据操作
type RbacStoreRepository struct {
	*BaseRepository[model.RbacStore]
}

// NewRbacStoreRepository 创建企业数据操作
func NewRbacStoreRepository(db *gorm.DB) *RbacStoreRepository {
	baseRepo := NewBaseRepository[model.RbacStore](db)
	return &RbacStoreRepository{
		BaseRepository: baseRepo,
	}
}

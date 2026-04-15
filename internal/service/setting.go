package service

import (
	"context"
	"encoding/json"
	"zero-backend/internal/apperror"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// SettingService 设置服务
type SettingService struct {
	repo        *repository.SettingRepository
	defaultRepo *repository.SettingDefaultRepository
}

// NewSettingService 创建设置服务
func NewSettingService(
	repo *repository.SettingRepository,
	defaultRepo *repository.SettingDefaultRepository,
) *SettingService {
	return &SettingService{
		repo:        repo,
		defaultRepo: defaultRepo,
	}
}

// FindList 获取设置列表
func (s *SettingService) FindList(ctx context.Context, req *dto.SettingListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.Setting{},
		Total: 0,
	}

	filter := &repository.SettingFilterField{
		SettingKey: req.SettingKey,
		StoreId:    req.StoreId,
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询设置数量失败")
	}

	if total == 0 {
		return result, nil
	}

	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询设置列表失败")
	}

	result.List = list

	return result, nil
}

// Create 创建设置
func (s *SettingService) Create(ctx context.Context, req *dto.SettingCreateRequest) error {
	if err := s.checkSettingKey(ctx, req.SettingKey, req.StoreId); err != nil {
		return err
	}

	item := &model.Setting{
		SettingKey:    req.SettingKey,
		SettingValues: req.SettingValues,
		Description:   req.Description,
		StoreId:       req.StoreId,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建设置失败")
	}

	return nil
}

// Update 更新设置
func (s *SettingService) Update(ctx context.Context, req *dto.SettingUpdateRequest) error {
	filter := &repository.SettingFilterField{
		Id:      req.ID,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询设置失败")
	}
	if item == nil || item.ID == 0 {
		return apperror.NewUserError("设置不存在或无权限访问")
	}

	if item.SettingKey != req.SettingKey {
		if err := s.checkSettingKey(ctx, req.SettingKey, req.StoreId); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"setting_key":    req.SettingKey,
		"setting_values": req.SettingValues,
		"description":    req.Description,
		"store_id":       req.StoreId,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新设置失败")
	}

	return nil
}

// Delete 删除设置
func (s *SettingService) Delete(ctx context.Context, req *dto.SettingDeleteRequest) error {
	filter := &repository.SettingFilterField{
		Id:      req.ID,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询设置失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除设置失败")
	}

	return nil
}

// checkSettingKey 检查设置key是否已存在
func (s *SettingService) checkSettingKey(ctx context.Context, key string, storeId uint32) error {
	filter := repository.SettingFilterField{SettingKey: key, StoreId: storeId}
	item, err := s.repo.FindOne(ctx, filter)

	if err != nil {
		return apperror.NewSystemError(err, "查询设置key失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("设置key已存在")
	}

	return nil
}

// GetSettingValue 获取并解析设置项
func (s *SettingService) GetSettingValue(ctx context.Context, key string, out interface{}) error {
	filter := repository.SettingFilterField{SettingKey: key}
	setting, err := s.repo.FindOne(ctx, filter)

	if err != nil {
		return apperror.NewSystemError(err, "查询设置失败")
	}

	var settingValues string
	if setting.ID > 0 {
		settingValues = setting.SettingValues
	} else {
		// 尝试从默认设置获取
		filter := &repository.SettingDefaultFilterField{SettingKey: key}
		defaultSetting, err := s.defaultRepo.FindOne(ctx, filter, repository.WithScopes(nil))

		if err != nil {
			return apperror.NewSystemError(err, "查询默认设置失败")
		}
		if defaultSetting.ID == 0 {
			return apperror.NewUserError("设置项不存在")
		}
		settingValues = defaultSetting.SettingValues
	}

	if err := json.Unmarshal([]byte(settingValues), out); err != nil {
		return apperror.NewSystemError(err, "解析设置值失败")
	}
	return nil
}

// FormConfigs 获取设置表单配置
func (s *SettingService) FormConfigs(ctx context.Context, req *dto.SettingFormConfigsRequest) ([]dto.SettingFormGroup, error) {
	configs := s.GetSettingFormConfigs()

	if req.OnlyPlatform {
		filtered := make([]dto.SettingFormGroup, 0)
		for _, config := range configs {
			if !config.OnlyPlatform {
				filtered = append(filtered, config)
			}
		}
		return filtered, nil
	}

	// 管理员用户返回全部配置
	return configs, nil
}

// QiniuToken 获取七牛上传token
func (s *SettingService) QiniuToken(ctx context.Context) (*dto.QiniuTokenResponse, error) {
	qiniu := &dto.QiniuConfig{}
	err := s.GetSettingValue(ctx, "qiniu", qiniu)
	if err != nil {
		return nil, err
	}

	putPolicy := storage.PutPolicy{
		Scope: qiniu.Bucket,
	}
	mac := qbox.NewMac(qiniu.AccessKey, qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	return &dto.QiniuTokenResponse{
		Token:     upToken,
		Domain:    qiniu.Domain,
		UploadUrl: "https://upload.qiniup.com",
	}, nil
}

// SettingDefaultService 默认设置服务
type SettingDefaultService struct {
	repo *repository.SettingDefaultRepository
}

// NewSettingDefaultService 创建默认设置服务
func NewSettingDefaultService(repo *repository.SettingDefaultRepository) *SettingDefaultService {
	return &SettingDefaultService{repo: repo}
}

// FindList 获取默认设置列表
func (s *SettingDefaultService) FindList(ctx context.Context, req *dto.SettingDefaultListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.SettingDefault{},
		Total: 0,
	}

	filter := &repository.SettingDefaultFilterField{
		SettingKey: req.SettingKey,
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询默认设置数量失败")
	}

	if total == 0 {
		return result, nil
	}

	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders, repository.WithScopes(nil))
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询默认设置列表失败")
	}

	result.List = list

	return result, nil
}

// Create 创建默认设置
func (s *SettingDefaultService) Create(ctx context.Context, req *dto.SettingDefaultCreateRequest) error {
	if err := s.checkSettingKey(ctx, req.SettingKey); err != nil {
		return err
	}

	item := &model.SettingDefault{
		SettingKey:    req.SettingKey,
		SettingValues: req.SettingValues,
		Description:   req.Description,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.NewSystemError(err, "创建默认设置失败")
	}

	return nil
}

// Update 更新默认设置
func (s *SettingDefaultService) Update(ctx context.Context, req *dto.SettingDefaultUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询默认设置失败")
	}

	if item.SettingKey != req.SettingKey {
		if err := s.checkSettingKey(ctx, req.SettingKey); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"setting_key":    req.SettingKey,
		"setting_values": req.SettingValues,
		"description":    req.Description,
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.NewSystemError(err, "更新默认设置失败")
	}

	return nil
}

// Delete 删除默认设置
func (s *SettingDefaultService) Delete(ctx context.Context, req *dto.SettingDefaultDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, repository.WithScopes(nil))
	if err != nil {
		return apperror.NewSystemError(err, "查询默认设置失败")
	}

	if item.ID == 0 {
		return apperror.NewUserError("找不到此记录")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除默认设置失败")
	}

	return nil
}

// checkSettingKey 检查默认设置key是否已存在
func (s *SettingDefaultService) checkSettingKey(ctx context.Context, key string) error {
	filter := &repository.SettingDefaultFilterField{SettingKey: key}
	item, err := s.repo.FindOne(ctx, filter, repository.WithScopes(nil))

	if err != nil {
		return apperror.NewSystemError(err, "查询默认设置key失败")
	}

	if item.ID > 0 {
		return apperror.NewUserError("默认设置key已存在")
	}

	return nil
}

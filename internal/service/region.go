package service

import (
	"context"
	"zero-backend/internal/apperror"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
)

// RegionService 通用服务
type RegionService struct {
	repo        *repository.RegionRepository
	settingServ *SettingService
}

// NewRegionService 创建通用服务
func NewRegionService(repo *repository.RegionRepository, settingServ *SettingService) *RegionService {
	return &RegionService{repo: repo, settingServ: settingServ}
}

// Regions 获取省市区数据(树形结构)
func (s *RegionService) Regions(ctx context.Context) ([]*model.Region, error) {
	// 获取所有地区数据
	regions, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询地区数据失败")
	}

	// 转换为RegionList并构建树形结构
	regionList := model.RegionList(regions)
	return regionList.Tree(), nil
}

package service

import (
	"context"
	"zero-backend/internal/apperror"
	"zero-backend/internal/constants"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
	"zero-backend/pkg/helper"

	"gorm.io/gorm"
)

// UserService 用户业务逻辑层
type UserService struct {
	db            *gorm.DB
	repo          *repository.UserRepository
	pointsLogRepo *repository.UserPointsLogRepository
}

// NewUserService 创建用户业务逻辑层
func NewUserService(db *gorm.DB, repo *repository.UserRepository, pointsLogRepo *repository.UserPointsLogRepository) *UserService {
	return &UserService{
		db:            db,
		repo:          repo,
		pointsLogRepo: pointsLogRepo,
	}
}

// List 获取用户列表
func (s *UserService) List(ctx context.Context, req *dto.UserListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.User{},
		Total: 0,
	}

	filter := &repository.UserFilterField{
		StoreId:  req.StoreId,
		Username: req.Username,
		Mobile:   req.Mobile,
		Status:   req.Status,
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	result.Total = total
	if total == 0 {
		return result, nil
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "id", Sort: "desc"},
	}

	list, err := s.repo.FindAll(ctx, filter, pagination, orders, repository.WithPreloads("Avatar"))
	if err != nil {
		return nil, err
	}

	result.List = list
	return result, nil
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest) error {
	if err := s.checkUsername(ctx, req.Username); err != nil {
		return err
	}

	// 密码加密
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return apperror.NewSystemError(err, "密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Mobile:   req.Mobile,
		NickName: req.NickName,
		AvatarId: req.AvatarId,
		Gender:   req.Gender,
		Status:   req.Status,
		StoreId:  req.StoreId,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return apperror.NewSystemError(err, "创建用户失败")
	}

	return nil
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, req *dto.UserUpdateRequest) error {
	filter := &repository.UserFilterField{
		Id:      req.Id,
		StoreId: req.StoreId,
	}
	user, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询用户失败")
	}

	if user.ID == 0 {
		return apperror.NewUserError("用户不存在或无权限访问")
	}

	if user.Username != req.Username {
		if err := s.checkUsername(ctx, req.Username); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"username":  req.Username,
		"mobile":    req.Mobile,
		"nick_name": req.NickName,
		"avatar_id": req.AvatarId,
		"gender":    req.Gender,
		"status":    req.Status,
		"store_id":  req.StoreId,
	}

	if req.Password != "" {
		hashedPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			return apperror.NewSystemError(err, "密码加密失败")
		}
		updateData["password"] = hashedPassword
	}

	if err := s.repo.Updates(ctx, user, updateData); err != nil {
		return apperror.NewSystemError(err, "更新用户失败")
	}

	return nil
}

// checkUsername 检查用户名是否已存在(同一企业内)
func (s *UserService) checkUsername(ctx context.Context, username string) error {
	filter := &repository.UserFilterField{Username: username}
	user, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询用户名失败")
	}

	if user.ID > 0 {
		return apperror.NewUserError("用户名已存在")
	}

	return nil
}

// GetPointsLogs 获取用户积分记录
func (s *UserService) GetPointsLogs(ctx context.Context, req *dto.UserPointsLogListRequest) (*dto.ListResult, error) {
	result := &dto.ListResult{
		List:  []*model.UserPointsLog{},
		Total: 0,
	}

	filter := &repository.UserPointsLogFilterField{
		StoreId: req.StoreId,
		UserId:  req.UserId,
	}

	total, err := s.pointsLogRepo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询积分记录总数失败")
	}

	result.Total = total
	if total == 0 {
		return result, nil
	}

	pagination := &repository.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	orders := repository.Orders{
		{Field: "created_at", Sort: "desc"},
	}

	list, err := s.pointsLogRepo.FindAll(ctx, filter, pagination, orders)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询积分记录失败")
	}

	result.List = list
	return result, nil
}

// ChangeUserPoints 变更用户积分
func (s *UserService) ChangeUserPoints(ctx context.Context, req *dto.UserPointsChangeRequest) error {
	// 参数校验
	if req.Points <= 0 {
		return apperror.NewUserError("积分变更值必须为正整数")
	}

	if !constants.PointsSourceType(req.SourceType).IsValid() {
		return apperror.NewUserError("无效的积分来源类型")
	}

	// 根据变更类型调整积分值
	points := req.Points
	if req.ChangeType == int8(constants.PointsChangeTypeReduce) {
		points = -points
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查询用户
		userFilter := &repository.UserFilterField{
			Id:      req.UserId,
			StoreId: req.StoreId,
		}
		user, err := s.repo.FindOne(ctx, userFilter, repository.WithTx[*repository.QueryConfig](tx))
		if err != nil {
			return apperror.NewSystemError(err, "查询用户失败")
		}
		if user.ID == 0 {
			return apperror.NewUserError("用户不存在")
		}

		// 更新用户积分
		updateData := map[string]any{
			"points": gorm.Expr("points + ?", points),
		}

		if err := s.repo.Updates(ctx, user, updateData, repository.WithTx[*repository.UpdateConfig](tx)); err != nil {
			return apperror.NewSystemError(err, "更新用户积分失败")
		}

		// 创建积分记录
		pointsLog := &model.UserPointsLog{
			UserId:     req.UserId,
			Points:     req.Points,
			ChangeType: req.ChangeType,
			SourceType: req.SourceType,
			SourceId:   req.SourceId,
			Remark:     req.Remark,
			StoreId:    req.StoreId,
		}

		if err := s.pointsLogRepo.Create(ctx, pointsLog, repository.WithTx[*repository.CreateConfig](tx)); err != nil {
			return apperror.NewSystemError(err, "创建积分记录失败")
		}

		return nil
	})
}

// Detail 获取用户详情
func (s *UserService) Detail(ctx context.Context, id uint32) (*model.User, error) {
	user, err := s.repo.FindOne(ctx, id)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询用户详情失败")
	}

	if user == nil || user.ID == 0 {
		return nil, apperror.NewUserError("用户不存在")
	}

	return user, nil
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, req *dto.UserDeleteRequest) error {
	filter := &repository.UserFilterField{
		Id:      req.Id,
		StoreId: req.StoreId,
	}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		return apperror.NewSystemError(err, "查询用户失败")
	}

	if item == nil || item.ID == 0 {
		return apperror.NewUserError("用户不存在或无权限访问")
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.NewSystemError(err, "删除用户失败")
	}
	return nil
}

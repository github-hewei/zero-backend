package user

import (
	"context"
	"errors"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"github.com/241x/zero-web/errcode"
	"gorm.io/gorm"
)

const PointsChangeTypeReduce = 2

// Service 用户业务逻辑
type Service struct {
	db            *gorm.DB
	repo          *Repository
	pointsLogRepo *PointsLogRepo
}

func NewService(db *gorm.DB, repo *Repository, pointsLogRepo *PointsLogRepo) *Service {
	return &Service{db: db, repo: repo, pointsLogRepo: pointsLogRepo}
}

func (s *Service) List(ctx context.Context, req *ListRequest) (*ListResult, error) {
	result := &ListResult{List: []*User{}, Total: 0}
	filter := &Filter{StoreId: req.StoreId, Username: req.Username, Mobile: req.Mobile, Status: req.Status}
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}
	result.Total = total
	if total == 0 {
		return result, nil
	}
	orders := baserepo.Orders{{Field: "id", Sort: "desc"}}
	list, err := s.repo.FindAll(ctx, filter, baserepo.NewPagination(req.Page, req.Limit), orders, baserepo.WithPreloads("Avatar"))
	if err != nil {
		return nil, err
	}
	result.List = list
	return result, nil
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) error {
	if err := s.checkUsername(ctx, req.Username); err != nil {
		return err
	}
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建用户失败"))
	}
	user := &User{Username: req.Username, Password: hashedPassword, Mobile: req.Mobile, NickName: req.NickName, AvatarId: req.AvatarId, Gender: req.Gender, Status: req.Status, StoreId: req.StoreId}
	if err := s.repo.Create(ctx, user); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建用户失败"))
	}
	return nil
}

func (s *Service) Update(ctx context.Context, req *UpdateRequest) error {
	filter := &Filter{Id: req.Id, StoreId: req.StoreId}
	user, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在或无权限访问"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新用户失败"))
	}
	if user.Username != req.Username {
		if err := s.checkUsername(ctx, req.Username); err != nil {
			return err
		}
	}
	updateData := map[string]any{
		"username": req.Username, "mobile": req.Mobile, "nick_name": req.NickName,
		"avatar_id": req.AvatarId, "gender": req.Gender, "status": req.Status, "store_id": req.StoreId,
	}
	if req.Password != "" {
		hashedPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新用户失败"))
		}
		updateData["password"] = hashedPassword
	}
	if err := s.repo.Updates(ctx, user, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新用户失败"))
	}
	return nil
}

func (s *Service) checkUsername(ctx context.Context, username string) error {
	filter := &Filter{Username: username}
	_, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("检查用户名失败"))
	}
	return apperror.New(errcode.Conflict, apperror.WithMsg("用户名已存在"))
}

func (s *Service) GetPointsLogs(ctx context.Context, req *PointsLogListRequest) (*ListResult, error) {
	result := &ListResult{List: []*PointsLog{}, Total: 0}
	filter := &PointsLogFilter{StoreId: req.StoreId, UserId: req.UserId}
	total, err := s.pointsLogRepo.Count(ctx, filter)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取积分记录失败"))
	}
	result.Total = total
	if total == 0 {
		return result, nil
	}
	orders := baserepo.Orders{{Field: "created_at", Sort: "desc"}}
	list, err := s.pointsLogRepo.FindAll(ctx, filter, baserepo.NewPagination(req.Page, req.Limit), orders)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取积分记录失败"))
	}
	result.List = list
	return result, nil
}

func (s *Service) ChangePoints(ctx context.Context, req *PointsChangeRequest) error {
	if req.Points <= 0 {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("积分变更值必须为正整数"))
	}
	if req.SourceType < 10 || req.SourceType > 30 {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("无效的积分来源类型"))
	}
	points := req.Points
	if req.ChangeType == PointsChangeTypeReduce {
		points = -points
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		userFilter := &Filter{Id: req.UserId, StoreId: req.StoreId}
		user, err := s.repo.FindOne(ctx, userFilter, baserepo.WithDB[*baserepo.QueryConfig](tx))
		if err != nil {
			if errors.Is(err, baserepo.ErrRecordNotFound) {
				return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
			}
			return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("变更用户积分失败"))
		}
		updateData := map[string]any{"points": gorm.Expr("points + ?", points)}
		if err := s.repo.Updates(ctx, user, updateData, baserepo.WithDB[*baserepo.UpdateConfig](tx)); err != nil {
			return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("变更用户积分失败"))
		}
		pointsLog := &PointsLog{UserId: req.UserId, Points: req.Points, ChangeType: req.ChangeType, SourceType: req.SourceType, SourceId: req.SourceId, Remark: req.Remark, StoreId: req.StoreId}
		if err := s.pointsLogRepo.Create(ctx, pointsLog, baserepo.WithDB[*baserepo.CreateConfig](tx)); err != nil {
			return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("变更用户积分失败"))
		}
		return nil
	})
}

func (s *Service) Detail(ctx context.Context, id uint32) (*User, error) {
	user, err := s.repo.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil, apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取用户详情失败"))
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, req *DeleteRequest) error {
	filter := &Filter{Id: req.Id, StoreId: req.StoreId}
	item, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在或无权限访问"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除用户失败"))
	}
	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除用户失败"))
	}
	return nil
}

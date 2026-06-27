package region

import (
	"context"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/errcode"
)

// Service 区划服务
type Service struct {
	repo *Repository
}

// NewService 创建区划服务
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Tree 获取区划树
func (s *Service) Tree(ctx context.Context) ([]*Region, error) {
	regions, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取地区数据失败"))
	}
	return List(regions).Tree(), nil
}

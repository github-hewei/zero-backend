package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zero-backend/internal/config"
	"zero-backend/internal/constants"
	"zero-backend/internal/dto"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/errcode"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/click"
)

// CaptchaService 验证码服务
type CaptchaService struct {
	rdb  *redis.Client
	cfg  config.CaptchaConfig
	capt click.Captcha
}

// NewCaptchaService 创建验证码服务实例
func NewCaptchaService(rdb *redis.Client, cfg config.CaptchaConfig) *CaptchaService {
	// 加载内嵌字体（方正盛世楷书简体）
	font, err := fzshengsksjw.GetFont()
	if err != nil {
		panic(fmt.Sprintf("captcha: 加载字体失败: %v", err))
	}

	// 加载内嵌背景图片
	bgImages, err := imagesv2.GetImages()
	if err != nil {
		panic(fmt.Sprintf("captcha: 加载背景图片失败: %v", err))
	}

	// 构建验证码实例
	builder := click.NewBuilder()
	builder.SetResources(
		click.WithFonts([]*truetype.Font{font}),
		click.WithBackgrounds(bgImages),
	)

	return &CaptchaService{
		rdb:  rdb,
		cfg:  cfg,
		capt: builder.Make(),
	}
}

// Generate 生成验证码
func (s *CaptchaService) Generate(ctx context.Context) (*dto.CaptchaGenerateResponse, error) {
	captData, err := s.capt.Generate()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	// 获取点选坐标数据（正确答案）
	dots := captData.GetData()
	if dots == nil {
		return nil, apperror.New(errcode.Internal, apperror.WithMsg("验证码数据为空"))
	}

	// 序列化正确答案
	dotsJSON, err := json.Marshal(dots)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	// 生成唯一 ID，存入 Redis
	captchaID := uuid.New().String()
	key := fmt.Sprintf("%s:%s", constants.RedisCaptchaKey, captchaID)

	if err := s.rdb.Set(ctx, key, dotsJSON, time.Duration(s.cfg.TTL)*time.Second).Err(); err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	// 生成 base64 图片
	masterImage, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码图片失败"))
	}

	thumbImage, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码图片失败"))
	}

	return &dto.CaptchaGenerateResponse{
		CaptchaID:   captchaID,
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
	}, nil
}

// Verify 校验验证码
func (s *CaptchaService) Verify(ctx context.Context, captchaID string, captchaCode string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisCaptchaKey, captchaID)

	// 从 Redis 取正确答案
	dotsJSON, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码已过期，请刷新重试"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证验证码失败"))
	}

	// 验证后立即删除（一次性消费，防重放）
	defer s.rdb.Del(ctx, key)

	// 反序列化正确答案
	var dots map[int]*click.Dot
	if err := json.Unmarshal(dotsJSON, &dots); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证验证码失败"))
	}

	// 反序列化前端提交的点击坐标
	var clickPoints []dto.CaptchaClickPoint
	if err := json.Unmarshal([]byte(captchaCode), &clickPoints); err != nil {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码数据格式错误"))
	}

	// 数量校验
	if len(clickPoints) != len(dots) {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码点选数量不正确"))
	}

	// 按顺序逐个校验坐标（容差 5px）
	for i, point := range clickPoints {
		dot, ok := dots[i]
		if !ok {
			return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码校验失败"))
		}
		if !click.Validate(point.X, point.Y, dot.X, dot.Y, dot.Width, dot.Height, 5) {
			return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码校验失败"))
		}
	}

	return nil
}

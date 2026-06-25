package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/errcode"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/click"
)

// Service 验证码服务
type Service struct {
	rdb    *redis.Client
	cfg    Config
	prefix string
	capt   click.Captcha
}

// NewService 创建验证码服务实例
func NewService(rdb *redis.Client, cfg Config, prefix string) *Service {
	font, err := fzshengsksjw.GetFont()
	if err != nil {
		panic(fmt.Sprintf("captcha: 加载字体失败: %v", err))
	}

	bgImages, err := imagesv2.GetImages()
	if err != nil {
		panic(fmt.Sprintf("captcha: 加载背景图片失败: %v", err))
	}

	builder := click.NewBuilder()
	builder.SetResources(
		click.WithFonts([]*truetype.Font{font}),
		click.WithBackgrounds(bgImages),
	)

	return &Service{
		rdb:    rdb,
		cfg:    cfg,
		prefix: prefix,
		capt:   builder.Make(),
	}
}

// Generate 生成验证码
func (s *Service) Generate(ctx context.Context) (*GenerateResponse, error) {
	captData, err := s.capt.Generate()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	dots := captData.GetData()
	if dots == nil {
		return nil, apperror.New(errcode.Internal, apperror.WithMsg("验证码数据为空"))
	}

	dotsJSON, err := json.Marshal(dots)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	captchaID := uuid.New().String()
	key := fmt.Sprintf("%s:%s", s.prefix, captchaID)

	if err := s.rdb.Set(ctx, key, dotsJSON, time.Duration(s.cfg.TTL)*time.Second).Err(); err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码失败"))
	}

	masterImage, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码图片失败"))
	}

	thumbImage, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("生成验证码图片失败"))
	}

	return &GenerateResponse{
		CaptchaID:   captchaID,
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
	}, nil
}

// Verify 校验验证码
func (s *Service) Verify(ctx context.Context, captchaID string, captchaCode string) error {
	key := fmt.Sprintf("%s:%s", s.prefix, captchaID)

	dotsJSON, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码已过期，请刷新重试"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证验证码失败"))
	}

	defer s.rdb.Del(ctx, key)

	var dots map[int]*click.Dot
	if err := json.Unmarshal(dotsJSON, &dots); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证验证码失败"))
	}

	var clickPoints []ClickPoint
	if err := json.Unmarshal([]byte(captchaCode), &clickPoints); err != nil {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码数据格式错误"))
	}

	if len(clickPoints) != len(dots) {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("验证码点选数量不正确"))
	}

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

package upload

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/errcode"
)

// Uploader 定义文件上传接口
type Uploader interface {
	Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (domain string, err error)
	Delete(ctx context.Context, filePath string) error
}

// NewUploader 创建上传器实例
func NewUploader(storageType string, ctx context.Context) (Uploader, error) {
	switch storageType {
	case "local":
		return NewLocalUploader(), nil
	case "qiniu":
		qiniuConfig := QiniuConfigFromCtx(ctx)
		return NewQiniuUploader(qiniuConfig), nil
	default:
		return nil, apperror.New(errcode.Internal, apperror.WithMsg("不支持的存储类型"))
	}
}

// LocalUploader 本地文件上传实现
type LocalUploader struct{}

// NewLocalUploader 创建本地文件上传器实例
func NewLocalUploader() *LocalUploader {
	return &LocalUploader{}
}

// Upload 上传文件到本地
func (u *LocalUploader) Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("打开文件失败"))
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(savePath), 0750); err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建目录失败"))
	}

	out, err := os.Create(savePath)
	if err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建文件失败"))
	}
	defer out.Close()

	if _, err = src.Seek(0, 0); err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("读取文件失败"))
	}

	if _, err = io.Copy(out, src); err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("写入文件失败"))
	}

	return "", nil
}

// Delete 删除本地文件
func (u *LocalUploader) Delete(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}

// qiniuConfigKey 上下文传递七牛云配置
type qiniuConfigKey struct{}

// WithQiniuConfig 注入七牛云配置
func WithQiniuConfig(ctx context.Context, cfg *QiniuConfig) context.Context {
	return context.WithValue(ctx, qiniuConfigKey{}, cfg)
}

// QiniuConfigFromCtx 从上下文读取七牛云配置
func QiniuConfigFromCtx(ctx context.Context) *QiniuConfig {
	v, ok := ctx.Value(qiniuConfigKey{}).(*QiniuConfig)
	if !ok {
		return nil
	}
	return v
}

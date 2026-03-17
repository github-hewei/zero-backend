package uploader

import (
	"context"
	"mime/multipart"
	"zero-backend/internal/apperror"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/dto"
)

// Uploader 定义文件上传接口
type Uploader interface {
	// Upload 执行文件上传
	Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (domain string, err error)

	// Delete 删除已上传文件(可选)
	Delete(ctx context.Context, filePath string) error
}

// NewUploader 创建上传器实例
func NewUploader(storageType string, ctx context.Context) (Uploader, error) {
	switch storageType {
	case "local":
		return NewLocalUploader(), nil
	case "qiniu":
		qiniuConfig, _ := ctx.Value(ctxkeys.QiniuConfigKey{}).(*dto.QiniuConfig)
		return NewQiniuUploader(qiniuConfig), nil
	default:
		return nil, apperror.NewUserError("不支持的存储类型")
	}
}

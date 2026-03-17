package uploader

import (
	"context"
	"mime/multipart"
	"strings"
	"zero-backend/internal/apperror"
	"zero-backend/internal/dto"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuUploader 七牛云文件上传实现
type QiniuUploader struct {
	config *dto.QiniuConfig
}

// NewQiniuUploader 创建七牛云上传实现
func NewQiniuUploader(config *dto.QiniuConfig) *QiniuUploader {
	return &QiniuUploader{config: config}
}

// getZone 获取区域
func (u *QiniuUploader) getZone(zone string) *storage.Zone {
	switch zone {
	case "z0":
		return &storage.ZoneHuadong
	case "z1":
		return &storage.ZoneHuabei
	case "z2":
		return &storage.ZoneHuanan
	case "na0":
		return &storage.ZoneBeimei
	case "as0":
		return &storage.ZoneXinjiapo
	default:
		return &storage.ZoneHuadong
	}
}

// Upload 上传文件到七牛云
func (u *QiniuUploader) Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (string, error) {
	// 初始化上传凭证
	putPolicy := storage.PutPolicy{
		Scope: u.config.Bucket,
	}
	mac := qbox.NewMac(u.config.AccessKey, u.config.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	// 配置上传参数
	cfg := storage.Config{
		Zone:          u.getZone(u.config.Zone),
		UseHTTPS:      true,
		UseCdnDomains: true,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 执行上传
	fileReader, err := file.Open()
	if err != nil {
		return "", apperror.NewSystemError(err, "打开上传文件失败")
	}
	defer fileReader.Close()

	savePath = strings.ReplaceAll(savePath, "\\", "/")
	err = formUploader.Put(ctx, &ret, upToken, savePath, fileReader, file.Size, nil)
	if err != nil {
		return "", apperror.NewSystemError(err, "七牛云上传失败")
	}

	return u.config.Domain, nil
}

// Delete 删除七牛云文件
func (u *QiniuUploader) Delete(ctx context.Context, filePath string) error {
	mac := qbox.NewMac(u.config.AccessKey, u.config.SecretKey)
	cfg := storage.Config{
		Zone: u.getZone(u.config.Zone),
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager.Delete(u.config.Bucket, filePath)
}

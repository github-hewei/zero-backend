package uploader

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-web/errcode"
)

// LocalUploader 本地文件上传实现
type LocalUploader struct{}

func NewLocalUploader() *LocalUploader {
	return &LocalUploader{}
}

// Upload 上传文件
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

// Delete 删除文件
func (u *LocalUploader) Delete(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}

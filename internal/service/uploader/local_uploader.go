package uploader

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"zero-backend/internal/errcode"
	"zero-backend/pkg/apperror"
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
		return "", apperror.Wrap(errcode.Internal, err)
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(savePath), 0750); err != nil {
		return "", apperror.Wrap(errcode.Internal, err)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return "", apperror.Wrap(errcode.Internal, err)
	}
	defer out.Close()

	if _, err = src.Seek(0, 0); err != nil {
		return "", apperror.Wrap(errcode.Internal, err)
	}

	if _, err = io.Copy(out, src); err != nil {
		return "", apperror.Wrap(errcode.Internal, err)
	}

	return "", nil
}

// Delete 删除文件
func (u *LocalUploader) Delete(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}

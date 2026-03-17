package dto

import "mime/multipart"

// UploadGroupCreateRequest 创建分组请求参数
type UploadGroupCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	ParentId uint32 `json:"parent_id"`
	Sort     uint32 `json:"sort"`
	StoreId  uint32 `json:"-"`
}

// UploadGroupUpdateRequest 更新分组请求参数
type UploadGroupUpdateRequest struct {
	ID       uint32 `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	ParentId uint32 `json:"parent_id"`
	Sort     uint32 `json:"sort"`
	StoreId  uint32 `json:"-"`
}

// UploadGroupDeleteRequest 删除分组请求参数
type UploadGroupDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// UploadFileListRequest 获取文件列表请求参数
type UploadFileListRequest struct {
	GroupId  string `json:"group_id"`
	FileType int8   `json:"file_type"`
	FileName string `json:"file_name"`
	Page     int    `json:"page" validate:"required,min=1"`
	Limit    int    `json:"limit" validate:"required,min=1,max=100"`
}

// UploadFileDeleteRequest 删除文件请求参数
type UploadFileDeleteRequest struct {
	ID uint32 `json:"id" validate:"required"`
}

// UploadFileRequest 文件上传请求参数
type UploadFileRequest struct {
	File    *multipart.FileHeader `json:"-"`        // 上传的文件
	GroupId uint32                `json:"group_id"` // 文件分组ID
}

// QiniuTokenResponse 七牛上传token
type QiniuTokenResponse struct {
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	UploadUrl string `json:"upload_url"`
}

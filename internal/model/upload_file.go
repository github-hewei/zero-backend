package model

// UploadFile 上传文件投影，仅包含用户模块需要的最小字段集合。
type UploadFile struct {
	ID       uint32 `json:"id" gorm:"primaryKey"`
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	Domain   string `json:"domain"`
}

func (UploadFile) TableName() string { return TableNamePrefix + "upload_file" }

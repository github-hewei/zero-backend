package repository

import (
	"zero-backend/internal/model"
	"zero-backend/pkg/helper"

	"gorm.io/gorm"
)

// UploadGroupFilterField 文件分组筛选字段
type UploadGroupFilterField struct {
	Name string
}

// Apply 应用筛选条件
func (f *UploadGroupFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.Name != "" {
		db = db.Where("name = ?", f.Name)
	}

	return db
}

// UploadGroupRepository 文件分组仓库
type UploadGroupRepository struct {
	*BaseRepository[model.UploadGroup]
}

// NewUploadGroupRepository 创建文件分组仓库
func NewUploadGroupRepository(db *gorm.DB) *UploadGroupRepository {
	baseRepo := NewBaseRepository[model.UploadGroup](db)
	return &UploadGroupRepository{
		BaseRepository: baseRepo,
	}
}

// UploadFileFilterField 文件查询过滤字段
type UploadFileFilterField struct {
	GroupId  string
	FileType int8
	FileName string
}

// Apply 应用筛选条件
func (f *UploadFileFilterField) Apply(db *gorm.DB) *gorm.DB {
	if f == nil {
		return db
	}

	if f.GroupId != "" && f.GroupId != "all" {
		db = db.Where("group_id = ?", f.GroupId)
	}

	if f.FileType > 0 {
		db = db.Where("file_type = ?", f.FileType)
	}

	if f.FileName != "" {
		db = db.Where("file_name like ?", helper.SafeLikeString(f.FileName)+"%")
	}

	return db
}

// UploadFileRepository 文件仓库
type UploadFileRepository struct {
	*BaseRepository[model.UploadFile]
}

// NewUploadFileRepository 创建文件仓库
func NewUploadFileRepository(db *gorm.DB) *UploadFileRepository {
	baseRepo := NewBaseRepository[model.UploadFile](db)
	return &UploadFileRepository{
		BaseRepository: baseRepo,
	}
}

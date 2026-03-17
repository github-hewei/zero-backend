package model

import "gorm.io/plugin/soft_delete"

type UploadGroup struct {
	ID        uint32 `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:30;not null;default:'';comment:分组名称"`
	ParentId  uint32 `json:"parent_id" gorm:"not null;default:0;comment:上级分组ID"`
	Sort      uint32 `json:"sort" gorm:"not null;default:0;comment:排序 ( 数字越小越靠前 ) "`
	StoreId   uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:store_id"`
	CreatedAt int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`

	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"not null;default:0;comment:删除时间"`
	Children  []*UploadGroup        `json:"children" gorm:"-"`
}

// TableName 指定数据表名称
func (m *UploadGroup) TableName() string {
	return TableNamePrefix + "upload_group"
}

// UploadGroupList 分组列表类型
type UploadGroupList []*UploadGroup

// Tree 转换为树形结构
func (list UploadGroupList) Tree() []*UploadGroup {
	// 创建ID到分组的映射
	groupMap := make(map[uint32]*UploadGroup)
	for _, group := range list {
		groupMap[group.ID] = group
	}

	// 构建树形结构
	var tree []*UploadGroup
	for _, group := range list {
		if group.ParentId == 0 {
			tree = append(tree, group)
		} else {
			if parent, ok := groupMap[group.ParentId]; ok {
				if parent.Children == nil {
					parent.Children = make([]*UploadGroup, 0)
				}
				parent.Children = append(parent.Children, group)
			}
		}
	}

	return tree
}

type UploadFile struct {
	ID         uint32 `json:"id" gorm:"primaryKey"`
	GroupId    uint32 `json:"group_id" gorm:"not null;default:0;comment:文件分组ID;index:group_id"`
	Channel    int8   `json:"channel" gorm:"type:tinyint;not null;default:10;comment:上传来源 ( 10后台 20客户端 ) "`
	Storage    string `json:"storage" gorm:"size:10;not null;default:'';comment:存储方式"`
	Domain     string `json:"domain" gorm:"size:255;not null;default:'';comment:存储域名"`
	FileType   int8   `json:"file_type" gorm:"type:tinyint;not null;default:10;comment:文件类型 ( 10图片 20附件 30视频 ) "`
	FileName   string `json:"file_name" gorm:"size:255;not null;default:'';comment:文件名称 ( 仅显示 ) "`
	FilePath   string `json:"file_path" gorm:"size:255;not null;default:'';comment:文件路径"`
	FileSize   uint32 `json:"file_size" gorm:"not null;default:0;comment:文件大小 ( 字节 ) "`
	FileExt    string `json:"file_ext" gorm:"size:20;not null;default:'';comment:文件扩展名"`
	Cover      string `json:"cover" gorm:"size:255;not null;default:'';comment:文件封面"`
	UploaderId uint32 `json:"uploader_id" gorm:"not null;default:0;comment:上传者用户ID"`
	StoreId    uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:store_id"`
	CreatedAt  int64  `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt  int64  `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`

	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"not null;default:0;comment:删除时间"`
}

// TableName 指定数据表名称
func (m *UploadFile) TableName() string {
	return TableNamePrefix + "upload_file"
}

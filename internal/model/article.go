package model

import "gorm.io/plugin/soft_delete"

// ArticleCategory 文章分类
type ArticleCategory struct {
	ID        uint32 `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:50;not null;default:'';comment:分类名称"`
	Status    int8   `json:"status" gorm:"type:tinyint;not null;default:1;comment:状态 ( 1显示 0隐藏 ) "`
	Sort      uint32 `json:"sort" gorm:"not null;default:0;comment:排序方式 ( 数字越小越靠前 ) "`
	StoreId   uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:store_id"`
	CreatedAt uint32 `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt uint32 `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`
}

// TableName 指定数据表名称
func (m *ArticleCategory) TableName() string {
	return TableNamePrefix + "article_category"
}

// Article 文章
type Article struct {
	ID           uint32 `json:"id" gorm:"primaryKey"`
	Title        string `json:"title" gorm:"size:300;not null;default:'';comment:文章标题"`
	ShowType     int8   `json:"show_type" gorm:"type:tinyint;not null;default:10;comment:列表显示方式 ( 10小图展示 20大图展示 ) "`
	CategoryId   uint32 `json:"category_id" gorm:"not null;default:0;comment:文章分类ID;index:category_id"`
	ImageId      uint32 `json:"image_id" gorm:"not null;default:0;comment:封面图ID"`
	Content      string `json:"content" gorm:"not null;comment:文章内容"`
	Sort         uint32 `json:"sort" gorm:"not null;default:0;comment:文章排序 ( 数字越小越靠前 ) "`
	Status       int8   `json:"status" gorm:"type:tinyint;not null;default:1;comment:文章状态 ( 0隐藏 1显示 ) "`
	VirtualViews uint32 `json:"virtual_views" gorm:"not null;default:0;comment:虚拟阅读量 ( 仅用作展示 ) "`
	ActualViews  uint32 `json:"actual_views" gorm:"not null;default:0;comment:实际阅读量"`
	StoreId      uint32 `json:"store_id" gorm:"not null;default:0;comment:企业ID;index:store_id"`
	CreatedAt    uint32 `json:"created_at" gorm:"not null;comment:创建时间;autoCreateTime"`
	UpdatedAt    uint32 `json:"updated_at" gorm:"not null;comment:更新时间;autoUpdateTime"`

	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"not null;default:0;comment:是否删除"`
	Image     *UploadFile           `json:"image" gorm:"foreignKey:ImageId"`
}

// TableName 指定数据表名称
func (m *Article) TableName() string {
	return TableNamePrefix + "article"
}

package dto

// ArticleCategoryListRequest 文章分类列表请求参数
type ArticleCategoryListRequest struct {
	Name    string `json:"name"`
	Status  int8   `json:"status" validate:"oneof=0 1"`
	StoreId uint32 `json:"store_id"`
	Page    int    `json:"page" validate:"required,min=1"`
	Limit   int    `json:"limit" validate:"required,min=1,max=100"`
}

// ArticleCategoryCreateRequest 创建文章分类请求参数
type ArticleCategoryCreateRequest struct {
	Name    string `json:"name" validate:"required,max=50"`
	Status  int8   `json:"status" validate:"required,oneof=0 1"`
	Sort    uint32 `json:"sort"`
	StoreId uint32 `json:"store_id"`
}

// ArticleCategoryUpdateRequest 更新文章分类请求参数
type ArticleCategoryUpdateRequest struct {
	Id      uint32 `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required,max=50"`
	Status  int8   `json:"status" validate:"required,oneof=0 1"`
	Sort    uint32 `json:"sort"`
	StoreId uint32 `json:"store_id"`
}

// ArticleCategoryDeleteRequest 删除文章分类请求参数
type ArticleCategoryDeleteRequest struct {
	Id uint32 `json:"id" validate:"required"`
}

// ArticleListRequest 文章列表请求参数
type ArticleListRequest struct {
	Title      string `json:"title"`
	CategoryId uint32 `json:"category_id"`
	Status     int8   `json:"status" validate:"oneof=0 1"`
	StoreId    uint32 `json:"store_id"`
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=100"`
}

// ArticleCreateRequest 创建文章请求参数
type ArticleCreateRequest struct {
	Title        string `json:"title" validate:"required,max=300"`
	ShowType     int8   `json:"show_type" validate:"required,oneof=10 20"`
	CategoryId   uint32 `json:"category_id" validate:"required"`
	ImageId      uint32 `json:"image_id"`
	Content      string `json:"content" validate:"required"`
	Sort         uint32 `json:"sort" validate:"required"`
	Status       int8   `json:"status" validate:"required,oneof=0 1"`
	VirtualViews uint32 `json:"virtual_views"`
	StoreId      uint32 `json:"store_id"`
}

// ArticleUpdateRequest 更新文章请求参数
type ArticleUpdateRequest struct {
	Id           uint32 `json:"id" validate:"required"`
	Title        string `json:"title" validate:"required,max=300"`
	ShowType     int8   `json:"show_type" validate:"required,oneof=10 20"`
	CategoryId   uint32 `json:"category_id" validate:"required"`
	ImageId      uint32 `json:"image_id"`
	Content      string `json:"content" validate:"required"`
	Sort         uint32 `json:"sort" validate:"required"`
	Status       int8   `json:"status" validate:"required,oneof=0 1"`
	VirtualViews uint32 `json:"virtual_views"`
	StoreId      uint32 `json:"store_id"`
}

// ArticleDeleteRequest 删除文章请求参数
type ArticleDeleteRequest struct {
	Id uint32 `json:"id" validate:"required"`
}

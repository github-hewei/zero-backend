package user

// ListResult 列表数据结构体
type ListResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Status   int8   `json:"status"`
	StoreId  uint32 `json:"store_id"`
	Page     int    `json:"page" validate:"required,min=1"`
	Limit    int    `json:"limit" validate:"required,min=1,max=100"`
}

// CreateRequest 创建用户请求参数
type CreateRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Mobile   string `json:"mobile" validate:"required,len=11"`
	NickName string `json:"nick_name" validate:"required,max=64"`
	AvatarId uint32 `json:"avatar_id"`
	Gender   int8   `json:"gender" validate:"required,oneof=0 1 2"`
	Status   int8   `json:"status" validate:"required,oneof=1 2"`
	StoreId  uint32 `json:"store_id"`
}

// UpdateRequest 更新用户请求参数
type UpdateRequest struct {
	Id       uint32 `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=5,max=32"`
	Password string `json:"password" validate:"min=6,max=32"`
	Mobile   string `json:"mobile" validate:"required,len=11"`
	NickName string `json:"nick_name" validate:"required,max=64"`
	AvatarId uint32 `json:"avatar_id"`
	Gender   int8   `json:"gender" validate:"required,oneof=0 1 2"`
	Status   int8   `json:"status" validate:"required,oneof=1 2"`
	StoreId  uint32 `json:"store_id"`
}

// DeleteRequest 删除用户请求参数
type DeleteRequest struct {
	Id      uint32 `json:"id" validate:"required"`
	StoreId uint32 `json:"store_id"`
}

// PointsLogListRequest 用户积分记录请求参数
type PointsLogListRequest struct {
	UserId     uint32 `json:"user_id" validate:"required"`
	StoreId    uint32 `json:"store_id"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	ChangeType int8   `json:"change_type" validate:"oneof=0 1 2"`
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=100"`
}

// PointsChangeRequest 用户积分变更请求参数
type PointsChangeRequest struct {
	UserId     uint32 `json:"user_id" validate:"required"`
	Points     int32  `json:"points" validate:"required"`
	ChangeType int8   `json:"change_type" validate:"required,oneof=1 2"`
	SourceType int8   `json:"source_type"`
	SourceId   string `json:"source_id" validate:"max=50"`
	Remark     string `json:"remark" validate:"max=255"`
	StoreId    uint32 `json:"store_id"`
}

// DetailRequest 用户详情请求参数
type DetailRequest struct {
	Id uint32 `json:"id" validate:"required"`
}

package dto

// UserListRequest 用户列表请求参数
type UserListRequest struct {
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Status   int8   `json:"status"`
	StoreId  uint32 `json:"store_id"`
	Page     int    `json:"page" validate:"required,min=1"`
	Limit    int    `json:"limit" validate:"required,min=1,max=100"`
}

// UserCreateRequest 创建用户请求参数
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Mobile   string `json:"mobile" validate:"required,len=11"`
	NickName string `json:"nick_name" validate:"required,max=64"`
	AvatarId uint32 `json:"avatar_id"`
	Gender   int8   `json:"gender" validate:"required,oneof=0 1 2"`
	Status   int8   `json:"status" validate:"required,oneof=1 2"`
	StoreId  uint32 `json:"store_id"`
}

// UserUpdateRequest 更新用户请求参数
type UserUpdateRequest struct {
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

// UserDeleteRequest 删除用户请求参数
type UserDeleteRequest struct {
	Id      uint32 `json:"id" validate:"required"`
	StoreId uint32 `json:"store_id"`
}

// UserPointsLogListRequest 用户积分记录请求参数
type UserPointsLogListRequest struct {
	UserId     uint32 `json:"user_id" validate:"required"`
	StoreId    uint32 `json:"store_id"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	ChangeType int8   `json:"change_type" validate:"oneof=0 1 2"`
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=100"`
}

// UserPointsChangeRequest 用户积分变更请求参数
type UserPointsChangeRequest struct {
	UserId     uint32 `json:"user_id" validate:"required"`
	Points     int32  `json:"points" validate:"required"`
	ChangeType int8   `json:"change_type" validate:"required,oneof=1 2"`
	SourceType int8   `json:"source_type"`
	SourceId   string `json:"source_id" validate:"max=50"`
	Remark     string `json:"remark" validate:"max=255"`
	StoreId    uint32 `json:"store_id"`
}

// UserDetailRequest 用户详情请求参数
type UserDetailRequest struct {
	Id uint32 `json:"id" validate:"required"`
}

// UserAddressListRequest 用户收货地址列表请求参数
type UserAddressListRequest struct {
	UserId  uint32 `json:"user_id" validate:"required"`
	StoreId uint32 `json:"store_id"`
	Name    string `json:"name" validate:"max=30"`
	Phone   string `json:"phone" validate:"max=20"`
	Page    int    `json:"page" validate:"required,min=1"`
	Limit   int    `json:"limit" validate:"required,min=1,max=100"`
}

// UserAddressCreateRequest 创建用户收货地址请求参数
type UserAddressCreateRequest struct {
	Name       string `json:"name" validate:"required,max=30"`
	Phone      string `json:"phone" validate:"required,max=20"`
	ProvinceId uint32 `json:"province_id" validate:"required"`
	CityId     uint32 `json:"city_id" validate:"required"`
	RegionId   uint32 `json:"region_id" validate:"required"`
	Detail     string `json:"detail" validate:"required,max=255"`
	UserId     uint32 `json:"user_id" validate:"required"`
	StoreId    uint32 `json:"store_id"`
}

// UserAddressUpdateRequest 更新用户收货地址请求参数
type UserAddressUpdateRequest struct {
	Id         uint32 `json:"id" validate:"required"`
	Name       string `json:"name" validate:"required,max=30"`
	Phone      string `json:"phone" validate:"required,max=20"`
	ProvinceId uint32 `json:"province_id" validate:"required"`
	CityId     uint32 `json:"city_id" validate:"required"`
	RegionId   uint32 `json:"region_id" validate:"required"`
	Detail     string `json:"detail" validate:"required,max=255"`
	UserId     uint32 `json:"user_id" validate:"required"`
	StoreId    uint32 `json:"store_id"`
}

// UserAddressDeleteRequest 删除用户收货地址请求参数
type UserAddressDeleteRequest struct {
	Id     uint32 `json:"id" validate:"required"`
	UserId uint32 `json:"user_id" validate:"required"`
}

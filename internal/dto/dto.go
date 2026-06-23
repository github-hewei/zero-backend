package dto

// ListResult 列表数据结构体
type ListResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

package base

// PageResult 通用分页响应结构。
type PageResult[T any] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
}

package vo

// PageInfo 分页信息
type PageInfo struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
}

// SearchParams 搜索参数
type SearchParams struct {
	Keyword  string `url:"keyword"`   // 搜索关键词
	Page     int    `url:"page"`      // 页码
	PageSize int    `url:"page_size"` // 每页数量
	Status   string `url:"status"`    // 状态筛选
	SortBy   string `url:"sort_by"`   // 排序字段
	Order    string `url:"order"`     // 排序方向 asc/desc
}

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

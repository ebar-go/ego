package pagination

// Paginator 分页器
type Paginator struct {
	// 总条数
	TotalCount int `json:"total"`
	// 每页行数
	Limit int `json:"per_page"`
	// 当前页数
	CurrentPage int `json:"current_page"`
	// 总页数
	TotalPages int `json:"total_pages"`
}



// GetOffset 获取偏移量
func (p *Paginator) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

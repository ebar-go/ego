package pagination

import "math"

const (
	defaultCurrentPage = 1
	defaultLimit = 10
)

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


func (p *Paginator) SetTotalCount(total int) {
	p.TotalCount = total
	p.TotalPages = int(math.Ceil(float64(p.TotalCount) / float64(p.Limit)))
}


// GetOffset 获取偏移量
func (p *Paginator) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

func New(currentPage, limit int) *Paginator {

	if currentPage <= 0 {
		currentPage = defaultCurrentPage
	}

	if limit <= 0 {
		limit = defaultLimit
	}

	return &Paginator{CurrentPage: currentPage, Limit: limit}
}
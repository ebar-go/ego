package library

import "math"

// Pagination 分页器
type Pagination struct {
	TotalCount int `json:"total"` // 总条数
	CurrentCount int `json:"count"` // 当前页的数据项数量
	Limit int `json:"per_page"` // 每页行数
	CurrentPage int `json:"current_page"`
	TotalPages int `json:"total_pages"`
	Link interface{} `json:"link"`
	Items interface{} `json:"-"`
}

const (
	defaultLimit = 10
	defaultCurrentPage = 1
)

// NewPagination 获取分页实例
func NewPagination(totalCount, currentPage, limit int) (Pagination) {
	pagination := Pagination{
		TotalCount: totalCount,
		CurrentPage: currentPage,
		Limit: limit,
	}

	if pagination.Limit <= 0 {
		pagination.Limit = defaultLimit
	}

	if pagination.CurrentPage <= 0 {
		pagination.CurrentPage = defaultCurrentPage
	}

	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalCount) / float64(pagination.Limit))) //page总数
	if pagination.CurrentPage < pagination.TotalPages {
		pagination.CurrentCount = pagination.Limit
	}else if pagination.CurrentPage > pagination.TotalPages {
		pagination.CurrentCount = 0
	}else {
		pagination.CurrentCount = pagination.TotalCount - (pagination.Limit * (pagination.CurrentPage - 1))
	}

	return pagination

}

// SetCurrentCount 设置当前页的数据项数量
func (p *Pagination) SetCurrentCount(currentCount int) {
	p.CurrentCount = currentCount
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

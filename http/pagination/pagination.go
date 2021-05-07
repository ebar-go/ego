package pagination

import (
	"math"
)

// Paginator 分页器
type Paginator struct {
	// 总条数
	TotalCount int `json:"total"`
	// 当前页的数据项数量
	CurrentCount int `json:"count"`
	// 每页行数
	Limit int `json:"per_page"`
	// 当前页数
	CurrentPage int `json:"current_page"`
	// 总页数
	TotalPages int `json:"total_pages"`
}

const (
	defaultLimit       = 10
	defaultCurrentPage = 1
)

func newPaginator(currentPage, limit int) Paginator {
	paginator := Paginator{
		CurrentPage: currentPage,
		Limit:       limit,
	}

	if paginator.Limit <= 0 {
		paginator.Limit = defaultLimit
	}

	if paginator.CurrentPage <= 0 {
		paginator.CurrentPage = defaultCurrentPage
	}

	return paginator
}

// NewPaginator return instance of Paginator
func NewPaginator(totalCount, currentPage, limit int) Paginator {
	paginator := newPaginator(currentPage, limit)
	paginator.TotalCount = totalCount

	paginator.TotalPages = int(math.Ceil(float64(totalCount) / float64(paginator.Limit)))
	if paginator.CurrentPage < paginator.TotalPages {
		paginator.CurrentCount = paginator.Limit
	} else if paginator.CurrentPage > paginator.TotalPages {
		paginator.CurrentCount = 0
	} else {
		paginator.CurrentCount = paginator.TotalCount - paginator.GetOffset()
	}

	return paginator
}

// NewArrayPaginator return instance of Paginator and current page result
func NewArrayPaginator(items []interface{}, currentPage, limit int) (paginator Paginator, result []interface{}) {
	paginator = newPaginator(currentPage, limit)
	paginator.TotalCount = len(items)
	paginator.TotalPages = int(math.Ceil(float64(paginator.TotalCount) / float64(paginator.Limit)))

	start := paginator.GetOffset()
	end := start + paginator.Limit
	if end > paginator.TotalCount {
		end = paginator.TotalCount
	}

	if start < end {
		paginator.CurrentCount = end - start
		result = items[start:end]
	} else {
		paginator.CurrentCount = 0
	}

	return
}

// GetOffset 获取偏移量
func (p *Paginator) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

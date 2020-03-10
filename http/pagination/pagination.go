package pagination

import (
	"github.com/ebar-go/ego/utils/number"
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

// Paginate 根据总条数、当前页数、每页行数分页
func Paginate(totalCount, currentPage, limit int) Paginator {
	pagination := Paginator{
		TotalCount:  totalCount,
		CurrentPage: currentPage,
		Limit:       limit,
	}

	if pagination.Limit <= 0 {
		pagination.Limit = defaultLimit
	}

	if pagination.CurrentPage <= 0 {
		pagination.CurrentPage = defaultCurrentPage
	}

	pagination.TotalPages = number.Div(totalCount, pagination.Limit) //page总数
	if pagination.CurrentPage < pagination.TotalPages {
		pagination.CurrentCount = pagination.Limit
	} else if pagination.CurrentPage > pagination.TotalPages {
		pagination.CurrentCount = 0
	} else {
		pagination.CurrentCount = pagination.TotalCount - pagination.GetOffset()
	}

	return pagination
}

// PaginateSlice 根据切片分页
func PaginateSlice(items []interface{}, currentPage, limit int) (paginate Paginator, result []interface{}) {
	paginate.CurrentPage = currentPage
	paginate.Limit = limit
	paginate.TotalCount = len(items)

	if paginate.Limit <= 0 {
		paginate.Limit = defaultLimit
	}

	if paginate.CurrentPage <= 0 {
		paginate.CurrentPage = defaultCurrentPage
	}

	paginate.TotalPages = int(math.Ceil(float64(paginate.TotalCount) / float64(paginate.Limit))) //page总数

	low := paginate.GetOffset()
	high := number.Min(paginate.TotalCount, low+paginate.Limit)

	if low < high {
		paginate.CurrentCount = high - low
		result = items[low:high]
	} else {
		paginate.CurrentCount = 0
	}

	return
}

// GetOffset 获取偏移量
func (p *Paginator) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

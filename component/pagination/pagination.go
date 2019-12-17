package pagination

import (
	"math"
	"github.com/ebar-go/ego/helper"
)

// Paginator 分页器
type Paginator struct {
	TotalCount int `json:"total"` // 总条数
	CurrentCount int `json:"count"` // 当前页的数据项数量
	Limit int `json:"per_page"` // 每页行数
	CurrentPage int `json:"current_page"`
	TotalPages int `json:"total_pages"`
	Link interface{} `json:"link"`
}

const (
	defaultLimit = 10
	defaultCurrentPage = 1
)

// Paginate 获取分页实例
func Paginate(totalCount, currentPage, limit int) Paginator {
	pagination := Paginator{
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

	pagination.TotalPages = helper.Div(totalCount, limit) //page总数
	if pagination.CurrentPage < pagination.TotalPages {
		pagination.CurrentCount = pagination.Limit
	}else if pagination.CurrentPage > pagination.TotalPages {
		pagination.CurrentCount = 0
	}else {
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
	high := helper.Min(paginate.TotalCount, low + paginate.Limit)

	if low < high {
		paginate.CurrentCount= high - low
		result = items[low:high]
	}else {
		paginate.CurrentCount = 0
	}

	return
}

// GetOffset 获取偏移量
func (p *Paginator) GetOffset() int {
	return (p.CurrentPage - 1) * p.Limit
}

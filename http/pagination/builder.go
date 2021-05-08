package pagination

import "math"


const (
	defaultLimit       = 10
	defaultCurrentPage = 1
)

// Builder 构造器
type Builder interface {
	// Build 生成Paginator实例
	Build() *Paginator
	// SetCurrentPage 设置当前页面，默认为1
	SetCurrentPage(currentPage int) Builder
	// SetLimit 设置每页行数，默认为10
	SetLimit(limit int) Builder
}

type paginatorBuilder struct {
	totalCount int
	currentPage int
	limit int
}

func (builder *paginatorBuilder) calculatePageNumber() int {
	return int(math.Ceil(float64(builder.totalCount) / float64(builder.limit)))
}

func (builder *paginatorBuilder) Build() *Paginator {
	paginator := &Paginator{
		TotalCount:  builder.totalCount,
		Limit:       builder.limit,
		CurrentPage: builder.currentPage,
		TotalPages:  builder.calculatePageNumber(),
	}
	return paginator
}

func (builder *paginatorBuilder) SetCurrentPage(currentPage int) Builder {
	if currentPage > 0 {
		builder.currentPage = currentPage
	}

	return builder
}

func (builder *paginatorBuilder) SetLimit(limit int) Builder {
	if limit > 0 {
		builder.limit = limit
	}
	return builder
}

// NewBuilder 返回Builder实例
func NewBuilder(total int) Builder {
	return &paginatorBuilder{totalCount: total, currentPage: defaultCurrentPage, limit: defaultLimit}
}


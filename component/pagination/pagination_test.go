package pagination

import (
	"testing"
	"fmt"
	"github.com/magiconair/properties/assert"
)

// TestNewPagination 测试分页
func TestPaginate(t *testing.T) {
	totalCount := 100
	currentPage := 1
	limit := 10
	pagination := Paginate(totalCount,currentPage,limit)

	assert.Equal(t, totalCount, pagination.TotalCount)
	assert.Equal(t, currentPage, pagination.CurrentPage)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, 10 , pagination.TotalPages)
	fmt.Println(pagination)

}

func TestNewPaginationWithSlice(t *testing.T) {
	currentPage := 2
	limit := 10
	items := []interface{}{1,2,3,4,5,6,7,8,9,10,11}
	pagination, result := PaginateSlice(items,currentPage,limit)
	assert.Equal(t, len(items), pagination.TotalCount)
	assert.Equal(t, currentPage, pagination.CurrentPage)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, 2 , pagination.TotalPages)
	assert.Equal(t, result, []interface{}{11})
}

package pagination

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

// TestNewPaginate 测试分页
func TestNewPaginate(t *testing.T) {
	totalCount := 100
	currentPage := 0
	limit := 0
	pagination := NewPaginator(totalCount, currentPage, limit)

	assert.Equal(t, 10, pagination.TotalPages)
	fmt.Println(pagination)

}

func TestNewArrayPaginator(t *testing.T) {
	currentPage := 2
	limit := 10
	items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	pagination, result := NewArrayPaginator(items, currentPage, limit)
	assert.Equal(t, len(items), pagination.TotalCount)
	assert.Equal(t, currentPage, pagination.CurrentPage)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, 2, pagination.TotalPages)
	assert.Equal(t, result, []interface{}{11})
}

func TestPaginator_GetOffset(t *testing.T) {
	pagination := NewPaginator(20, 0, 0)
	assert.Equal(t, 0, pagination.GetOffset())
}
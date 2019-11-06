package library

import (
	"testing"
	"fmt"
	"github.com/magiconair/properties/assert"
)

// TestNewPagination 测试分页
func TestNewPagination(t *testing.T) {
	totalCount := 100
	currentPage := 1
	limit := 10
	pagination := NewPagination(totalCount,currentPage,limit)

	assert.Equal(t, totalCount, pagination.TotalCount)
	assert.Equal(t, currentPage, pagination.CurrentPage)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, 10 , pagination.TotalPages)
	fmt.Println(pagination)


}

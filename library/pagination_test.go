package library

import (
	"testing"
	"github.com/ebar-go/ego/test"
	"fmt"
)

// TestNewPagination 测试分页
func TestNewPagination(t *testing.T) {
	totalCount := 100
	currentPage := 1
	limit := 10
	pagination := NewPagination(totalCount,currentPage,limit)

	test.AssertEqual(t, totalCount, pagination.TotalCount)
	test.AssertEqual(t, currentPage, pagination.CurrentPage)
	test.AssertEqual(t, limit, pagination.Limit)
	test.AssertEqual(t, 10 , pagination.TotalPages)
	fmt.Println(pagination)


}

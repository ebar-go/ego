package pagination

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewBuilder 测试分页
func TestNewBuilder(t *testing.T) {
	paginator := NewBuilder(100).Build()
	assert.Equal(t, 1, paginator.CurrentPage)
	assert.Equal(t, 10, paginator.Limit)
	assert.Equal(t, 10, paginator.TotalPages)
}

func TestPaginator_GetOffset(t *testing.T) {
	paginator := NewBuilder(20).SetCurrentPage(3).SetLimit(4).Build()
	assert.Equal(t, 8, paginator.GetOffset())
}

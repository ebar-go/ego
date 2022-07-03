package pagination

import (
	"github.com/stretchr/testify/assert"
	"testing"
)



func TestPaginator_GetOffset(t *testing.T) {
	paginator := New(3, 4)
	paginator.SetTotalCount(20)
	assert.Equal(t, 8, paginator.GetOffset())
}

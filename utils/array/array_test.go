package array

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func getIntItems() []int {
	return []int{1,2,3,4}
}

func getStringItems() []string {
	return []string{"1","2","3","4"}
}

func TestInt(t *testing.T)  {
	s := getIntItems()
	a := Int(s)
	assert.Equal(t, a.items, s)
	assert.Equal(t, a.Length(), len(s))
	assert.Equal(t, a.Implode(","), "1,2,3,4")
	assert.Equal(t, a.ToString(), getStringItems())

	a.Push(2)
	assert.Equal(t, a.Length(), len(s)+1)
	assert.Equal(t, a.Unique(), s)
}

func TestString(t *testing.T)  {
	s := getStringItems()
	a := String(s)
	assert.Equal(t, a.items, s)
	assert.Equal(t, a.Length(), len(s))
	assert.Equal(t, a.Implode(","), "1,2,3,4")
	assert.Equal(t, a.ToInt(), getIntItems())

	a.Push("2")
	assert.Equal(t, a.Length(), len(s)+1)
	assert.Equal(t, a.Unique(), s)
}
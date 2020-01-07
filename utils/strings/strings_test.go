package strings

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

type Item struct {
	Data interface{}
	Expect string
}

var items = []Item{
	{
		Data:   []int{1,2,3},
		Expect: "1,2,3",
	},
	{
		Data:   []string{"hello","world"},
		Expect: "hello,world",
	},
	{
		Data:   []interface{}{"test", 1},
		Expect: "test,1",
	},
}

func TestImplode(t *testing.T) {
	for _, v := range items {
		got := Implode(",", v.Data)
		assert.Equal(t, v.Expect, got)
	}
}

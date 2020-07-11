package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Obj struct {
	Name string `json:"name"`
}

func TestJson(t *testing.T) {
	o := Obj{
		Name: "test",
	}

	s, err := Encode(o)
	assert.Nil(t, err)

	var d Obj
	assert.Nil(t, Decode([]byte(s), &d))
	assert.Equal(t, o, d)

}

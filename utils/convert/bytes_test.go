package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString2Byte(t *testing.T) {
	s := "foo"
	p := String2Byte(s)
	s2 := Byte2String(p)
	assert.Equal(t, s, s2)
}

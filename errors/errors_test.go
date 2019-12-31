package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(200105001, "something wrong" )
	assert.Error(t, err)
}

func TestError_Error(t *testing.T) {
	err := New(200105001, "something wrong" )
	assert.Error(t, err)
	fmt.Println(err.Error())
}

func TestParse(t *testing.T) {
	err := New(200105001, "something wrong" )
	assert.Error(t, err)

	errParse := Parse(err.Error())
	assert.Error(t, errParse)

}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("aha")
	assert.Error(t, err)
	pe := Parse(err.Error())
	assert.Equal(t, pe.Code, 500)
}
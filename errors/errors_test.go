package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(200105001, "something wrong")
	assert.Error(t, err)
}

func TestError_Error(t *testing.T) {
	err := New(200105001, "something wrong")
	assert.Error(t, err)
	fmt.Println(err.Error())
}


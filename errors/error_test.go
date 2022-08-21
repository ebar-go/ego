package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithMessage(t *testing.T) {
	assert.Nil(t, WithMessage(nil, "some message"))
	assert.NotNil(t, WithMessage(errors.New("some error"), "some message"))
}

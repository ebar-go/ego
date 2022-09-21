package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestWithMessage(t *testing.T) {
	assert.Nil(t, WithMessage(nil, "some message"))
	assert.NotNil(t, WithMessage(errors.New("some error"), "some message"))
}

func TestNew(t *testing.T) {
	log.Printf("err: %v\n", New(unknown, "New"))
	log.Printf("err: %v\n", WithMessage(errors.New("test"), "WithMessage"))
	log.Printf("err: %v\n", Convert(errors.New("Convert")))
	log.Printf("err: %v\n", NotFound("NotFound"))
}

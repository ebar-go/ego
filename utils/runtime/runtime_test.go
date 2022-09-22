package runtime

import (
	"errors"
	"testing"
)

func TestLogPanic(t *testing.T) {
	logPanic(errors.New("some panic"))
}

func TestHandleCrash(t *testing.T) {
	defer HandleCrash()
	panic(errors.New("some panic"))
}

func TestSetReallyCrash(t *testing.T) {
	defer HandleCrash()
	SetReallyCrash(false)
	panic(errors.New("some panic"))

}

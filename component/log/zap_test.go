package log

import (
	"testing"
)

func TestZap(t *testing.T) {

	Info("Info", Context{
		"hello": "world",
	})
	Debug("Debug", Context{
		"hello": "world",
	})
	Error("Error", Context{
		"hello": "world",
	})

}

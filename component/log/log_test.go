package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	logger := New("/tmp/app.log", true, nil)
	logger.Debug("test", Context{"hello":"world"})

}

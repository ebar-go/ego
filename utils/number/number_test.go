package number

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFloatValue_Int(t *testing.T) {
	f := FloatValue(1.0)
	assert.Equal(t, f.Int(), 1)
}

func TestFloatValue_Round(t *testing.T) {
	f := FloatValue(1.5)
	assert.Equal(t, f.Round(), 2)
}

package app

import (
	"github.com/ebar-go/ego/log"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, Config().ServicePort, 8080)
	assert.Equal(t, Config().LogPath, "/tmp")

	Config().ServicePort = 8081
	assert.Equal(t, Config().ServicePort, 8081)
}

func TestLogManager(t *testing.T) {
	LogManager().App().Debug("test", log.Context{"hello": "world"})
}

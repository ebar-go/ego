package app

import (
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/event"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContainer(t *testing.T) {
	assert.NotNil(t, NewContainer())
}

func TestConfig(t *testing.T) {
	assert.Equal(t, Config().ServicePort, 8080)
	assert.Equal(t, Config().LogPath, "/tmp")

	Config().ServicePort = 8081
	assert.Equal(t, Config().ServicePort, 8081)
}

func TestLogManager(t *testing.T) {
	LogManager().App().Debug("test", log.Context{"hello": "world"})
}

func TestTask(t *testing.T) {
	assert.IsType(t, Task(), &cron.Cron{})
}

func TestJwt(t *testing.T) {
	assert.IsType(t, Jwt(), &auth.JwtAuth{})
}

func TestEventDispatcher(t *testing.T) {
	assert.IsType(t, EventDispatcher(), event.NewDispatcher())
}

func TestMns(t *testing.T) {
	assert.NotNil(t, Mns())
}

func TestMysql(t *testing.T) {
	assert.NotNil(t, Mysql())
}

func TestRedis(t *testing.T) {
	assert.NotNil(t, Redis())
}

func TestWebSocket(t *testing.T) {
	assert.NotNil(t, WebSocket())
}

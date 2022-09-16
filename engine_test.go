package ego

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNamedEngine(t *testing.T) {
	name := "test"
	engine := NewNamedEngine(name)
	assert.NotNil(t, engine)
	assert.Equal(t, name, engine.name)
}

func TestNamedEngine_Run(t *testing.T) {
	engine := NewNamedEngine("test")
	assert.NotNil(t, engine)
	engine.Run()
}

func TestEngine_WithComponent(t *testing.T) {
	engine := buildEngine().WithComponent(component.NewCache())
	assert.NotNil(t, engine)
}

func TestEngine_WithServer(t *testing.T) {
	engine := buildEngine().WithServer(server.NewHTTPServer(":8080"))
	assert.NotNil(t, engine)
}

func TestEngine_Run(t *testing.T) {
	engine := buildEngine()
	assert.NotNil(t, engine)
	engine.Run()
}

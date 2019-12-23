package http

import (
	"fmt"
	"github.com/ebar-go/ego/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	app := NewServer()
	assert.NotNil(t, app)
}

func TestServer_Start(t *testing.T) {
	app := NewServer()

	// 添加路由
	app.Router.GET("/test", func(context *gin.Context) {
		fmt.Println("hello,world")
	})

	assert.Nil(t, app.Start(config.Instance.ServicePort))
}

func TestMain(m *testing.M) {

	m.Run()
}

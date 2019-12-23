package http

import (
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	assert.NotNil(t, server)
}

func TestServer_Start(t *testing.T) {
	server := NewServer()

	// 添加路由
	server.Router.GET("/test", func(context *gin.Context) {
		fmt.Println("hello,world")
	})

	err := server.Start(app.Config().ServicePort)
	assert.Nil(t, err)
}

func TestMain(m *testing.M) {

	m.Run()
}

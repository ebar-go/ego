package http

import (
	"fmt"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	app := NewServer()
	assert.Nil(t, app.Start())
}


func TestServer_Init(t *testing.T) {
	app := NewServer()
	app.SetName("test")
	app.SetJwtKey([]byte("jwt_key"))
	app.SetLogPath("/tmp/log")
	assert.Nil(t, app.Start())

	// 添加路由
	app.Router.Use(middleware.JWT)
	app.Router.GET("/test", func(context *gin.Context) {
		fmt.Println("hello,world")
	})

	assert.Nil(t, app.Start())
}

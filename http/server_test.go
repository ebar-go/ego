package http

import (
	"fmt"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_Init(t *testing.T) {
	server := &Server {
		Address : "127.0.0.1", // 可以读取apollo地址
		Port:8088,
		LogPath:"/tmp/log",
	}
	assert.Nil(t, server.Init())
}

func TestServer_GetCompleteHost(t *testing.T) {
	server := &Server {
		Address : "127.0.0.1", // 可以读取apollo地址
		Port:8088,
		LogPath:"/tmp/log",
	}
	assert.Equal(t, fmt.Sprintf("%s:%d", server.Address, server.Port), server.GetCompleteHost())
}

func TestServer_Start(t *testing.T) {
	server := &Server{
		Address : "127.0.0.1", // 可以读取apollo地址
		Port:8088,
		LogPath:"/tmp/log",
		JwtKey:[]byte("jwt_key"),
	}
	assert.Nil(t, server.Init())

	// 添加路由
	server.Router.Use(middleware.JWT)
	server.Router.GET("/test", func(context *gin.Context) {
		fmt.Println("hello,world")
	})

	assert.Nil(t, server.Start())
}

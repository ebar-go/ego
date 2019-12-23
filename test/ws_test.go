package test

import (
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/http"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/ws"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Listen(t *testing.T) {
	server := http.NewServer()

	// 添加路由
	server.Router.GET("/test", func(ctx *gin.Context) {
		response.Success(ctx, "hello")
	})

	server.Router.GET("/ws", func(ctx *gin.Context) {
		conn, err := ws.GetUpgradeConnection(ctx.Writer, ctx.Request)
		if err != nil {
			response.Error(ctx, 500, err.Error())
			return
		}
		fmt.Println(1)

		client := ws.NewClient(conn, func(ctx *ws.Context) string {
			return helper.GetTimeStr() + ":" + ctx.GetMessage()
		})

		client.CloseHandler = func() {
			fmt.Println("closed")
		}

		app.WebSocket().RegisterClient(client)
		go client.Listen()

	})

	go app.WebSocket().Start()

	err := server.Start(app.Config().ServicePort)
	assert.Nil(t, err)
}

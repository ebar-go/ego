package ego

import (
	"fmt"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/ws"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)



func TestWebsocketServer(t *testing.T) {
	httpServer := HttpServer()
	httpServer.Port = 9001
	httpServer.Router.Use(middleware.Recover)
	webSocket := WebsocketServer()

	httpServer.Router.GET("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})

	// 注册
	httpServer.Router.GET("/ws", func(ctx *gin.Context) {
		// get ws conn
		conn, err := ws.NewConnection(ctx.Writer, ctx.Request)
		if err != nil {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}

		_ = conn.Send([]byte(fmt.Sprintf("welcome:%s", conn.ID)))
		// 接收数据并处理
		conn.Handle(func(message []byte) {
			// 给其他人发信息
			content := fmt.Sprintf("%s talk:%s", conn.ID, string(message))
			webSocket.Broadcast(egu.Str2Byte(content), conn)

		})

		webSocket.Register(conn)
	})

	// 给客户端发送数据
	httpServer.Router.GET("/send", func(ctx *gin.Context) {
		id := ctx.Query("id")
		conn := webSocket.GetConnection(id)
		if conn == nil {
			panic(errors.NotFound("no connection"))
		}

		_ = conn.Send([]byte("haha"))
		response.WrapContext(ctx).Success(nil)
	})
	// 给客户端广播数据
	httpServer.Router.GET("/broadcast", func(ctx *gin.Context) {
		webSocket.Broadcast([]byte("hello,world"), nil)
		response.WrapContext(ctx).Success(nil)
	})

	// 模拟剔除在线用户
	httpServer.Router.GET("/delete", func(ctx *gin.Context) {
		id := ctx.Query("id")
		conn := webSocket.GetConnection(id)
		if conn == nil {
			panic(errors.NotFound("no connection"))
		}

		_ = conn.Send([]byte("haha"))
		webSocket.Unregister(id)
		response.WrapContext(ctx).Success(nil)
	})

	webSocket.Start()

	egu.FatalError("StartHttpServer", httpServer.Start())
}

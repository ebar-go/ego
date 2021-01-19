package ego

import (
	"github.com/ebar-go/ego/errors"
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
	webSocket := WebsocketServer()

	httpServer.Router.GET("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})
	httpServer.Router.GET("/ws", func(ctx *gin.Context) {
		// get ws conn
		conn, err := ws.NewConnection(ctx.Writer, ctx.Request)
		if err != nil {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}

		conn.Handle(func(message []byte) {
			_ = conn.Send(append([]byte("receive:"), message...))

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

		_ = conn.Send([]byte("send to user"))
		response.WrapContext(ctx).Success(nil)
	})
	// 给客户端发送数据
	httpServer.Router.GET("/broadcast", func(ctx *gin.Context) {
		webSocket.Broadcast([]byte("hello,world"), nil)
		response.WrapContext(ctx).Success(nil)
	})

	webSocket.Start()

	egu.FatalError("StartHttpServer", httpServer.Start())
}

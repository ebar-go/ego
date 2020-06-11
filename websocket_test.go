package ego

import (
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils/secure"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestWebsocketServer(t *testing.T) {
	s := HttpServer()
	ws := WebsocketServer()

	s.Router.GET("/check", func(context *gin.Context) {
		response.WrapContext(context).Success("hello")
	})
	s.Router.GET("/ws", func(ctx *gin.Context) {
		// get websocket conn
		conn, err := ws.UpgradeConn(ctx.Writer, ctx.Request)
		if err != nil {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}

		ws.Register(conn, func(message []byte){
			if string(message) == "broadcast" {// 广播
				ws.Broadcast([]byte("hello,welcome"), nil)
				return
			}
			ws.Send(message, conn)

		})
	})

	go ws.Start()

	secure.FatalError("StartHttpServer", s.Start())
}

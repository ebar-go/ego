package ego

import (
	"github.com/ebar-go/ego/ws"
)


// Websocket return ws websocketServer instance
func WebsocketServer() ws.Server {
	return ws.New()
}

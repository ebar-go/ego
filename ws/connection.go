package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

// GetUpgradeConnection get web socket connection
func GetUpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	respHeader := http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}}
	return u.Upgrade(w, r, respHeader)
}

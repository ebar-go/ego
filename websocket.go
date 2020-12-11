package ego

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// WsServer
type WsServer interface {
	Register(conn *WebsocketConn, handler Handler)
	Unregister(id string)
	Start()
	Broadcast(message []byte, ignore *WebsocketConn)
	Send(message []byte, c *WebsocketConn)
	UpgradeConn(w http.ResponseWriter, r *http.Request) (*WebsocketConn, error)
	GetOnline()(int)
	//GetOnlines()(*websocketServer)
	IsRegist(uuid string) (int)
	GetConnection(uuid string)(*WebsocketConn)
}

// Handler define message processor
type Handler func(message []byte)

// Websocket return ws websocketServer instance
func WebsocketServer() *websocketServer {
	return &websocketServer{
		connections: make(map[string]*WebsocketConn),
		register:    make(chan *WebsocketConn),
		unregister:  make(chan *WebsocketConn),
	}
}

// websocketServer implement ws websocketServer interface
type websocketServer struct {
	connections map[string]*WebsocketConn
	register    chan *WebsocketConn
	unregister  chan *WebsocketConn
}

// Register register conn
func (websocketServer *websocketServer) Register(conn *WebsocketConn, handler Handler) {
	conn.handler = handler
	go conn.listen()
	websocketServer.register <- conn
}

// UnRegister delete websocket connection
func (websocketServer *websocketServer) Unregister(id string) {
	conn, ok := websocketServer.connections[id]
	if ok {
		websocketServer.unregister <- conn
	}
}

// Start
func (websocketServer *websocketServer) Start() {
	for {
		select {
		case conn := <-websocketServer.register:
			websocketServer.connections[conn.ID] = conn

			conn.websocketServer = websocketServer
		case conn := <-websocketServer.unregister:
			if _, ok := websocketServer.connections[conn.ID]; ok {
				delete(websocketServer.connections, conn.ID)
			}
		}
	}
}

// Broadcast push message to all connection, except ignore connection
func (websocketServer *websocketServer) Broadcast(message []byte, ignore *WebsocketConn) {
	for id, conn := range websocketServer.connections {
		if ignore == nil || ignore.ID != id {
			websocketServer.Send(message, conn)
		}
	}
}

// Send push message to client
func (websocketServer *websocketServer) Send(message []byte, c *WebsocketConn) {
	_ = c.conn.WriteMessage(websocket.TextMessage, message)
}

var u = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

// WebsocketConn return web socket connection
func (websocketServer *websocketServer) UpgradeConn(w http.ResponseWriter, r *http.Request) (*WebsocketConn, error) {
	respHeader := http.Header{"Sec-WebSocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")}}
	conn, err := u.Upgrade(w, r, respHeader)
	if err != nil {
		return nil, err
	}

	return &WebsocketConn{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}, nil
}

// websocketConn include socket conn
type WebsocketConn struct {
	// unique id
	ID string

	// socket connection
	conn *websocket.Conn

	// process handler
	handler Handler

	// connection websocketServer
	websocketServer *websocketServer
}

func (c *WebsocketConn) GetID() string {
	return c.ID
}

// close
func (c *WebsocketConn) close() {
	c.websocketServer.Unregister(c.ID)
	_ = c.conn.Close()
}

// Listen listen connection
func (c *WebsocketConn) listen() {
	defer func() {
		c.close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.handler(message)
	}
}

func (websocketServer *websocketServer)GetOnline()(int){
	return len(websocketServer.connections)
}

func (websocketServer *websocketServer)IsRegist(uuid string) (int)  {
	res := 0
	if _,ok := websocketServer.connections[uuid];ok{
		res = 1
	}
	return res
}

func (websocketServer *websocketServer)GetConnection(uuid string)(*WebsocketConn){
	return websocketServer.connections[uuid]
}


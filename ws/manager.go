package ws

// IWebSocketManager ws管理器接口
type IWebSocketManager interface {
	Start()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	Broadcast(message string)
	Send(message []byte, ignore *Client)

}

// ClientManager is a websocket manager
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var manager = DefaultManager()

// GetManager
func GetManager() IWebSocketManager {
	return manager
}

// Manager define a ws server manager
func DefaultManager() IWebSocketManager {
	return &ClientManager{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (manager *ClientManager) RegisterClient(client *Client) {

	manager.register <- client
}

func (manager *ClientManager) UnregisterClient(client *Client) {
	client.OnClose()
	manager.unregister <- client
}

// Broadcast 官博
func (manager *ClientManager) Broadcast(message string) {
	manager.broadcast <- []byte(message)
}

// Start is to start a ws server
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register: // new connection
			manager.clients[conn] = true
			conn.OnOpen()

		case conn := <-manager.unregister: // close connection

			if _, ok := manager.clients[conn]; ok {
				close(conn.Send)
				delete(manager.clients, conn)
			}

		case message := <-manager.broadcast: // Broadcast
			for conn := range manager.clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

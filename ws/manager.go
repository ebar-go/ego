package ws

// Manager ws管理器接口
type Manager interface {
	// register client
	RegisterClient(client *Client)

	// unregister client
	UnregisterClient(client *Client)

	// broadcast message
	Broadcast(message string)

	// send message
	Send(message []byte, ignore *Client)

	Start()
}

// manager is a websocket manager
type manager struct {
	clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

// NewManager define a ws server manager
func NewManager() Manager {
	return &manager{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// RegisterClient 注册客户端
func (manager *manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnregisterClient 注销客户端
func (manager *manager) UnregisterClient(client *Client) {
	manager.Unregister <- client
}

// Start
func (manager *manager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.clients[client.ID] = client
			client.OnOpen()

			// emm..如果不这么做，没有想到更好的办法在defer时删除manager的client
			client.manager = manager
		case client := <-manager.Unregister:
			if _, ok := manager.clients[client.ID]; ok {
				delete(manager.clients, client.ID)
			}
		}
	}
}

// Broadcast 广播消息
func (manager *manager) Broadcast(message string) {
	for _, conn := range manager.clients {
		conn.SendMessage([]byte(message))
	}
}

// Send send message to clients and ignore someone
func (manager *manager) Send(message []byte, ignore *Client) {
	for id, conn := range manager.clients {
		if id != ignore.ID {
			conn.SendMessage(message)
		}
	}
}

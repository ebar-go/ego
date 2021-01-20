package ws

// Server
type Server interface {
	Register(conn *Connection)
	Unregister(id string)
	Start()
	Broadcast(message []byte, ignore *Connection)
	GetConnection(id string) *Connection
}

type server struct {
	connections map[string]*Connection
	register    chan *Connection
	unregister  chan *Connection
}

// New
func New() Server {
	return &server{
		connections: make(map[string]*Connection),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
	}
}


// Register register conn
func (srv *server) Register(conn *Connection) {
	go conn.listen(srv.unregister)
	srv.register <- conn
}

// UnRegister delete ws connection
func (srv *server) Unregister(id string) {
	conn, ok := srv.connections[id]

	if ok {
		conn.close(srv.unregister)
	}
}

// Start
func (srv *server) Start() {
	go func() {
		for {
			select {
			case conn := <-srv.register:
				srv.connections[conn.ID] = conn
			case conn := <-srv.unregister:
				if _, ok := srv.connections[conn.ID]; ok {
					delete(srv.connections, conn.ID)
				}
			}
		}
	}()

}

// Broadcast push message to all connection, except ignore connection
func (srv *server) Broadcast(message []byte, ignore *Connection) {
	for id, conn := range srv.connections {
		if ignore == nil || ignore.ID != id {
			_ = conn.Send(message)
		}
	}
}

// 获取连接，may be nil
func (srv *server) GetConnection(id string) *Connection{
	return srv.connections[id]
}
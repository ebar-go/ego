package server

type HTTPServer struct {
}

func (s HTTPServer) Serve(stop <-chan struct{}) error {
	//TODO implement me
	panic("implement me")
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{}
}

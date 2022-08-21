package server

type Server interface {
	Serve(stop <-chan struct{})
}

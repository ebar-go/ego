package protocol

import (
	"fmt"
	"net"
)

// Schema defines server protocol
type Schema struct {
	Protocol string
	Bind     string
}

// Address returns address of schema.
func (s Schema) Address() string {
	return fmt.Sprintf("%s://%s", s.Protocol, s.Bind)
}

// HostPort returns the host and port of the schema.
func (s Schema) HostPort() (string, string, error) {
	return net.SplitHostPort(s.Bind)
}

func NewSchema(protocol string, bind string) Schema {
	return Schema{Protocol: protocol, Bind: bind}
}

func NewHttpSchema(bind string) Schema {
	return NewSchema(HTTP, bind)
}
func NewGRPCSchema(bind string) Schema {
	return NewSchema(GRPC, bind)
}

func NewWSSchema(bind string) Schema {
	return NewSchema(WS, bind)
}

func NewTCPSchema(bind string) Schema {
	return NewSchema(TCP, bind)
}

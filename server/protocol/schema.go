package protocol

type Schema struct {
	Protocol string
	Bind     string
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

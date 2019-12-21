package client

import (
	"github.com/ebar-go/ego/http/client/impl"
	"github.com/ebar-go/ego/http/client/request"
)

// Client Http客户端
type Client interface {
	Execute(request request.IRequest) (string, error)
	NewRequest(param request.Param) request.IRequest
}

func NewHttpClient() impl.HttpClient {
	return impl.NewHttpClient()
}

func NewFastHttpClient() impl.FastHttpClient {
	return impl.NewFastHttpClient()
}

func NewKongClient() impl.KongClient {
	return impl.KongClient{}
}

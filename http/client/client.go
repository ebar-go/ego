package client

import (
	"github.com/ebar-go/ego/http/client/impl"
	"github.com/ebar-go/ego/http/client/request"
)

// IClient Http客户端
type IClient interface {
	Execute(request request.IRequest) (string, error)
	NewRequest(param request.Param) request.IRequest
}

func NewHttpClient() IClient {
	return impl.NewHttpClient()
}

func NewFastHttpClient() IClient {
	return impl.NewFastHttpClient()
}

func NewKongClient() IClient {
	return impl.KongClient{}
}

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

// NewHttpClient 官方http客户端
func NewHttpClient() impl.HttpClient {
	return impl.NewHttpClient()
}

// NewFastHttpClient fastHttp客户端
func NewFastHttpClient() impl.FastHttpClient {
	return impl.NewFastHttpClient()
}

// NewKongClient Kong网关客户端
func NewKongClient() impl.KongClient {
	return impl.KongClient{}
}

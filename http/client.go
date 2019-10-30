package http

import (
	"net/http"
	"time"
	"net"
	"io"
	"github.com/ebar-go/ego/http/middleware"
	"fmt"
)

// Client http客户端
type Client struct {
	Timeout time.Duration
	Transport *http.Transport
	initialize bool
	clientPool *http.Client
}

const(
	DefaultMaxIdleConns = 100
	DefaultMaxIdleConnsPerHost = 100
)

type KongRequest struct {
	Iss string // 签名
	Secret string // 秘钥
	ReferServiceName string
	ReferRequestHost string
	GatewayTrace string
	XServiceUser string
	TokenExpireTime int // jwt过滤地址
	Address string // kong网关地址
}

const(
	HttpSchema = "http://"
)

// NewRequest
func (kong *KongRequest) NewRequest(method , uri string, body io.Reader) *http.Request {
	url := HttpSchema + kong.Address + uri
	request, _ := http.NewRequest(method , url, body)

	jwtToken , _ := middleware.GetEncodeToken(kong.Iss, kong.Secret, kong.TokenExpireTime)

	request.Header.Add("Accept-Encoding", "charset=UTF-8")
	request.Header.Add("Refer-Service-Name", kong.ReferServiceName)
	request.Header.Add("Refer-Request-Host", kong.ReferRequestHost)
	request.Header.Add("X-Service-User", kong.XServiceUser)
	request.Header.Add(middleware.JwtTokenHeader, fmt.Sprintf("%s %s", middleware.JwtTokenMethod,jwtToken))
	fmt.Println(url, request.Header)

	return request
}

// DefaultClient 默认http客户端
func DefaultClient() *Client {
	return  &Client{
		Timeout: time.Duration(3) * time.Second,
		Transport: &http.Transport{ // 配置连接池
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        DefaultMaxIdleConns,
			MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(90) * time.Second,
		},
	}
}

// NewRequest
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}


// GetInstance 获取http请求示例
func (client *Client) GetInstance() *http.Client{
	 if !client.initialize {
	 	if client.Timeout == 0 {
	 		client.Timeout = time.Duration(2) * time.Second
		}

	 	httpClient := &http.Client{
			Timeout: client.Timeout,
		}

	 	if client.Transport != nil {
	 		httpClient.Transport = client.Transport
		}

	 	client.clientPool = httpClient

	 	client.initialize = true
	 }

	 return client.clientPool
}



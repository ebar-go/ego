package http

import (
	"net/http"
	"time"
)

// Client http客户端
type Client struct {
	Timeout time.Duration
	Transport *http.Transport
	initialize bool
	clientPool *http.Client
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



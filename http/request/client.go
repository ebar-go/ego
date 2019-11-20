package request

import (
	"time"
	"net/http"
	"net"
	"io"
	"io/ioutil"
	"github.com/pkg/errors"
	"strings"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/library"
)

const(
	DefaultMaxIdleConns = 100
	DefaultMaxIdleConnsPerHost = 100
)

// Client http客户端
type Client struct {
	Timeout time.Duration
	Transport *http.Transport
	clientPool *http.Client
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

// complete 补全默认值
func (client *Client) complete()  {
	if client.clientPool == nil {
		httpClient := &http.Client{
			Timeout: client.Timeout,
		}

		if client.Transport != nil {
			httpClient.Transport = client.Transport
		}

		client.clientPool = httpClient
	}

}

// New 实例化request
func New(method, url string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(url, constant.HttpSchema) {
		url = constant.HttpSchema + url
	}

	library.Debug(url)
	return http.NewRequest(method, url, body)
}

// Do 执行http请求
func (client *Client) Do(req *http.Request) (*http.Response, error){
	client.complete()
	return client.clientPool.Do(req)
}

// 将response序列化
func  StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("没有响应数据")
	}

	if response.StatusCode != 200 {
		return "", errors.New("非200的上游返回")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.WithMessage(err, "读取响应数据失败:")
	}

	// 关闭响应
	defer func() {
		response.Body.Close()
	}()

	return string(data), nil
}

package component

import (
	"errors"
	"fmt"
	"github.com/ebar-go/ego/rebuild/utils/serializer"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	EmptyResponse = errors.New("empty response")
)

// Curl simple wrapper of http.Client
type Curl struct {
	Named
	httpClient   *http.Client
	bufferLenght int
}

// Send return Response by http.Request
func (c *Curl) Send(request *http.Request) (serializer.Serializer, error) {
	resp, err := c.httpClient.Do(request)
	// close response
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	return c.ReadResponse(resp)
}

func (c *Curl) ReadResponse(resp *http.Response) (serializer.Serializer, error) {
	if resp == nil {
		return nil, EmptyResponse
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is:%d", resp.StatusCode)
	}

	length := int(resp.ContentLength)
	if length == 0 {
		length = c.bufferLenght
	}
	buffer := serializer.NewBuffer(length)

	_, err := io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return buffer, nil
}

func NewCurl() *Curl {
	return &Curl{
		Named:        Named("curl"),
		bufferLenght: 512,
		httpClient: &http.Client{
			Transport: &http.Transport{ // 配置连接池
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout: 3 * time.Second,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Second,
		},
	}
}

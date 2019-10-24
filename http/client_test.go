package http

import (
	"testing"
	"git.epetbar.com/go-package/ego/test"
	"fmt"
	"net/http"
	"net"
	"time"
)

// TestClient_GetInstance 测试http实例化
func TestClient_GetInstance(t *testing.T) {
	client := Client{
		Timeout: time.Duration(3) * time.Second,
		Transport: &http.Transport{ // 配置连接池
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     time.Duration(90) * time.Second,
		},
	}

	httpClient := client.GetInstance()
	test.AssertNotNil(t, httpClient)

	response, err := httpClient.Get("http://localhost:8088")
	fmt.Println(err)
	test.AssertNil(t, err)
	test.AssertNotNil(t, response)

	fmt.Println(StringifyResponse(response))

}

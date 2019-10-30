package http

import (
	"testing"
	"github.com/ebar-go/ego/test"
	"fmt"
	"net/http"
	"net"
	"time"
	"github.com/ebar-go/ego/library"
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

func TestKongRequest_NewRequest(t *testing.T) {
	kongRequest := KongRequest{
		Iss:"common-openapi",
		Secret:"WUcLklcyETkhz7ktThMniw6AFseNbrJ6",
		Address: "47.110.77.180:8000",
	}

	request := kongRequest.NewRequest("GET", "/gott-wms/v1/basicInformation/warehouse/list?ware_nos=163", nil)

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

	resp, err := httpClient.Do(request)
	fmt.Println(resp, err)
	str, _ :=StringifyResponse(resp)
	library.Debug(str)
}

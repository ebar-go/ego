package request

import (
	"io"
	"net/http"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/constant"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/pkg/errors"
	"github.com/ebar-go/ego/component/trace"
)

type Kong struct {
	Iss string // 签名
	Secret string // 秘钥
	ReferServiceName string
	ReferRequestHost string
	GatewayTrace string
	XServiceUser string
	TokenExpireTime int // jwt过滤地址
	Address string // kong网关地址
}

// NewRequest
func (kong *Kong) NewRequest(router, method , uri string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s%s", kong.Address, router, uri)
	request, err := New(method , url, body)
	if err != nil {
		return nil, err
	}


	// 生成kong的token
	jwtToken , _ := middleware.GetEncodeToken(kong.Iss, kong.Secret, kong.TokenExpireTime)

	request.Header.Add("Accept-Encoding", "charset=UTF-8")
	request.Header.Add("Refer-Service-Name", kong.ReferServiceName)
	request.Header.Add("Refer-Request-Host", kong.ReferRequestHost)
	request.Header.Add("X-Service-User", kong.XServiceUser)
	request.Header.Add(constant.JwtTokenHeader, fmt.Sprintf("%s %s", constant.JwtTokenMethod,jwtToken))

	return request, nil
}

func (kong *Kong) NewFastHttpRequest(router, method, uri string) *fasthttp.Request{
	url := fmt.Sprintf("%s/%s%s", kong.Address, router, uri)

	req := fasthttp.AcquireRequest()

	req.Header.SetContentType("application/json")
	req.Header.SetMethod(method)

	req.SetRequestURI(url)

	// 生成kong的token
	jwtToken , _ := middleware.GetEncodeToken(kong.Iss, kong.Secret, kong.TokenExpireTime)

	req.Header.Add("Accept-Encoding", "charset=UTF-8")
	req.Header.Add("Refer-Service-Name", kong.ReferServiceName)
	req.Header.Add("Refer-Request-Host", kong.ReferRequestHost)
	req.Header.Add("X-Service-User", kong.XServiceUser)
	req.Header.Add(constant.GatewayTrace, trace.GetTraceId())
	req.Header.Add(constant.JwtTokenHeader, fmt.Sprintf("%s %s", constant.JwtTokenMethod,jwtToken))

	return req
}

// Send 发送http请求，得到响应
func Send(request *fasthttp.Request) (string, error) {
	defer fasthttp.ReleaseRequest(request) // 用完需要释放资源

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := fasthttp.Do(request, resp); err != nil {
		return "", errors.WithMessage(err, "请求失败")
	}

	return string(resp.Body()), nil
}

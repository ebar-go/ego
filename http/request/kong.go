package request

import (
	"io"
	"net/http"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/http/constant"
	"fmt"
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

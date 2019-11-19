package request

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/http/constant"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"time"
)

type Kong struct {
	Iss string // 签名
	Secret string // 秘钥
	ReferServiceName string
	ReferRequestHost string
	GatewayTrace string
	XServiceUser string
	TokenExpireTime int // jwt过滤时间
	Address string // kong网关地址
}

// GetEncodeToken 获取加密的token
func (kong Kong) GenerateToken() (string, error) {
	if kong.TokenExpireTime == 0 {
		kong.TokenExpireTime = constant.JwtExpiredTime
	}
	now := time.Now().Unix()
	exp := now + int64(kong.TokenExpireTime)
	claim := jwt.MapClaims{
		"iss":       kong.Iss,
		"iat":      now,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claim)
	tokenStr,err  := token.SignedString([]byte(kong.Secret))
	return tokenStr, err
}

// NewRequest
func (kong *Kong) NewRequest(router, method , uri string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s%s", kong.Address, router, uri)
	request, err := New(method , url, body)
	if err != nil {
		return nil, err
	}


	// 生成kong的token
	jwtToken , _ := kong.GenerateToken()

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
	jwtToken , _ := kong.GenerateToken()

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

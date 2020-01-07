package impl

import (
	"fmt"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/http/client/request"
	"time"
	"github.com/dgrijalva/jwt-go"
)

// Kong kong客户端
type KongClient struct {
	Iss              string // 签名
	Secret           string // 秘钥
	ReferServiceName string
	ReferRequestHost string
	GatewayTrace     string
	XServiceUser     string
	TokenExpireTime  int    // jwt过滤时间
	Address          string // kong网关地址
}

// GetCompleteUrl 获取完整的地址
func (kong KongClient) GetCompleteUrl(router, uri string) string {
	return fmt.Sprintf("%s/%s%s", kong.Address, router, uri)
}

// GetEncodeToken 获取加密的token
func (kong KongClient) GenerateToken() (string, error) {
	if kong.TokenExpireTime == 0 {
		kong.TokenExpireTime = 3600
	}
	now := time.Now().Unix()
	exp := now + int64(kong.TokenExpireTime)
	claim := jwt.MapClaims{
		"iss": kong.Iss,
		"iat": now,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(kong.Secret))
	return tokenStr, err
}

// NewRequest
func (kong KongClient) NewRequest(param request.Param) request.IRequest {
	// 生成kong的token
	jwtToken, _ := kong.GenerateToken()
	param.AddHeader("Accept-Encoding", "charset=UTF-8")
	param.AddHeader("Refer-Service-Name", kong.ReferServiceName)
	param.AddHeader("Refer-Request-Host", kong.ReferRequestHost)
	param.AddHeader("X-Service-User", kong.XServiceUser)
	param.AddHeader("gateway-trace", trace.GetTraceId())
	param.AddHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

	return NewFastHttpClient().NewRequest(param)
}

// Send 发送http请求，得到响应
func (kong KongClient) Execute(req request.IRequest) (string, error) {

	return NewFastHttpClient().Execute(req)
}

package middleware

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

var jwtAuth = new(JwtAuth)

// JwtAuth jwt
type JwtAuth struct {
	SigningKey []byte
}

// SetJwtSigningKey 设置jwt的秘钥
func SetJwtSigningKey(key []byte) {
	jwtAuth.SigningKey = key
}

var TokenNotExist = errors.New("token not exist")
var TokenValidateFailed = errors.New("token validate failed")

// CreateToken 生成一个token
func (jwtAuth JwtAuth) CreateToken(claimsCreator func() jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsCreator())
	return token.SignedString(jwtAuth.SigningKey)
}

// parseToken 解析Token
func (jwtAuth JwtAuth) parseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, TokenValidateFailed
	}

	return tokenClaims.Claims, nil
}

// GetCurrentClaims 获取解析jwt后的信息
func GetCurrentClaims(ctx *gin.Context) interface{} {
	claims, exist := ctx.Get(constant.JwtClaimsKey)
	if !exist {
		return nil
	}

	return claims
}

// validateToken 验证token
func (jwtAuth JwtAuth) validateToken(ctx *gin.Context) error {
	// 获取token
	tokenStr := ctx.GetHeader(constant.JwtTokenHeader)
	kv := strings.Split(tokenStr, " ")
	if len(kv) != 2 || kv[0] != constant.JwtTokenMethod {
		return TokenNotExist
	}

	claims, err := jwtAuth.parseToken(kv[1])
	if err != nil {
		return err
	}

	// token存入context
	ctx.Set(constant.JwtClaimsKey, claims)
	return nil
}

// JWT gin的jwt中间件
func JWT(ctx *gin.Context) {

	// 解析token
	if err := jwtAuth.validateToken(ctx); err != nil {
		response.Error(ctx, constant.StatusUnauthorized, err.Error())

		ctx.Abort()
		return
	}

	ctx.Next()
}

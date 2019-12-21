package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/container"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	ClaimsKey = "jwt_claims"
)

// Jwt json web token
type Jwt interface {
	// 解析token
	ParseToken(token string) (jwt.Claims, error)

	// 创建token
	CreateToken(claimsCreator func() jwt.Claims) (string, error)
}

func NewJwt(signKey []byte) Jwt {
	return &JwtAuth{SignKey: signKey}
}

// JwtAuth jwt
type JwtAuth struct {
	SignKey []byte
}

var (
	TokenNotExist       = errors.New("token not exist")
	TokenValidateFailed = errors.New("token validate failed")
)

// CreateToken 生成一个token
func (jwtAuth JwtAuth) CreateToken(claimsCreator func() jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsCreator())
	return token.SignedString(jwtAuth.SignKey)
}

// ParseToken 解析Token
func (jwtAuth JwtAuth) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.SignKey, nil
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
	claims, exist := ctx.Get(ClaimsKey)
	if !exist {
		return nil
	}

	return claims
}

// validateToken 验证token
func validateToken(ctx *gin.Context) error {
	// 获取token
	tokenStr := ctx.GetHeader("Authorization")
	kv := strings.Split(tokenStr, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return TokenNotExist
	}

	return container.App.Invoke(func(jwtAuth Jwt) error {
		claims, err := jwtAuth.ParseToken(kv[1])
		if err != nil {
			return err
		}

		// token存入context
		ctx.Set(ClaimsKey, claims)
		return nil
	})
}

// JWT gin的jwt中间件
func JWT(ctx *gin.Context) {

	// 解析token
	if err := validateToken(ctx); err != nil {
		response.Error(ctx, 401, err.Error())

		ctx.Abort()
		return
	}

	ctx.Next()
}

package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	ClaimsKey = "jwt_claims"
)

var (
	TokenNotExist = errors.New("token not exist")
)

// GetCurrentClaims 获取解析jwt后的信息
func GetCurrentClaims(ctx *gin.Context) jwt.MapClaims {
	claims, exist := ctx.Get(ClaimsKey)
	if !exist || claims == nil {
		return nil
	}

	return claims.(jwt.MapClaims)
}

// validateToken 验证token
func validateToken(ctx *gin.Context) error {
	// 获取token
	tokenStr := ctx.GetHeader("Authorization")
	kv := strings.Split(tokenStr, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return TokenNotExist
	}

	claims, err := auth.New(config.Server().JwtSignKey).ParseToken(kv[1])
	if err != nil {
		return err
	}

	// token存入context
	ctx.Set(ClaimsKey, claims)
	return nil
}

// JWT gin的jwt中间件
func JWT(ctx *gin.Context) {
	// 解析token
	if err := validateToken(ctx); err != nil {
		response.WrapContext(ctx).Error(401, err.Error())

		ctx.Abort()
		return
	}

	ctx.Next()
}

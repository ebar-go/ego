package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	TokenNotExist = errors.New("token not exist")
)

// validateToken 验证token
func validateToken(ctx *gin.Context, claims jwt.Claims) error {
	// 获取token
	tokenStr := ctx.GetHeader("Authorization")
	kv := strings.Split(tokenStr, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return TokenNotExist
	}

	// parse claims
	if err := app.Jwt().ParseTokenWithClaims(kv[1], claims); err != nil {
		return err
	}

	// 令牌信息存入context
	ctx.Set("claims", claims)
	return nil
}

// JWT gin的jwt中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析token
		if err := validateToken(ctx, new(jwt.StandardClaims)); err != nil {
			response.WrapContext(ctx).Error(401, err.Error())

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

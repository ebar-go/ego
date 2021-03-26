package middleware

import (
	"fmt"
	"github.com/ebar-go/ego/component/auth"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWT gin的jwt中间件
func JWT(jwtAuth auth.Jwt) gin.HandlerFunc {
	validator := func(ctx *gin.Context) error {
		token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			return fmt.Errorf("invalid token")
		}

		claims, err := jwtAuth.ParseToken(token)
		if err != nil {
			return fmt.Errorf("parse token: %v", err)
		}

		// 令牌信息存入context
		ctx.Set("claims", claims)

		return nil
	}

	return func(ctx *gin.Context) {
		// 解析token
		if err := validator(ctx); err != nil {
			response.WrapContext(ctx).Error(401, err.Error())

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

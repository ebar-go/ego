package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
	"github.com/ebar-go/ego/http"
	"strings"
	"errors"
	"os"
)

var jwtSecret  []byte

// Claims Jwt自定义结构体
type Claims struct {
	User User
	jwt.StandardClaims
}

// User 用户数据
type User struct {
	UID int `json:"uid"`
	ACID string `json:"acId"`
	Name string `json:"name"`
}

var TokenNotExist = errors.New("token not exist")
var TokenValidateFailed = errors.New("token validate failed")
var TokenExpired = errors.New("token expired")

const(
	JwtTokenMethod = "Bearer"
	JwtTokenHeader = "Authorization"
)

// SetJwtSecret 设置jwt的秘钥
func SetJwtSecret(secret []byte)  {
	jwtSecret = secret
}


// ParseToken 解析Token
func ParseToken(token string) (*Claims, error) {
	if string(jwtSecret) == "" {
		jwtSecret = []byte(os.Getenv("JWT_KEY"))
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	fmt.Println(1, tokenClaims, err)

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}


// JWT gin的jwt中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		// 获取token
		token := c.GetHeader(JwtTokenHeader)

		kv := strings.Split(token, " ")
		if len(kv) != 2 || kv[0] != JwtTokenMethod {
			err = TokenNotExist
		}else {
			token = kv[1]
			claims, err := ParseToken(token)
			if err != nil {
				err = TokenValidateFailed
			} else if time.Now().Unix() > claims.ExpiresAt {
				err = TokenExpired
			}
		}

		if err != nil {
			response := http.DefaultResponse(c)
			response.StatusCode = 401
			response.Message = err.Error()
			c.JSON(401, response)

			c.Abort()
			return
		}

		c.Next()
	}
}
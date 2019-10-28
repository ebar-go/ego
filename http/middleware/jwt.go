package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
	"strings"
	"errors"
	"os"
	"github.com/ebar-go/ego/http/response"
)

var jwtSecret []byte

// Claims Jwt自定义结构体
type Claims struct {
	User User
	jwt.StandardClaims
}

// User 用户数据
type User struct {
	UID  int    `json:"uid"`
	ACID string `json:"acId"`
	Name string `json:"name"`
}

var TokenNotExist = errors.New("token not exist")
var TokenValidateFailed = errors.New("token validate failed")
var TokenExpired = errors.New("token expired")

const (
	JwtTokenMethod = "Bearer"
	JwtTokenHeader = "Authorization"
)

// SetJwtSecret 设置jwt的秘钥
func SetJwtSecret(secret []byte) {
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

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// JWT gin的jwt中间件
func JWT(c *gin.Context) {
	var errRes error

	// 获取token
	token := c.GetHeader(JwtTokenHeader)

	kv := strings.Split(token, " ")
	if len(kv) != 2 || kv[0] != JwtTokenMethod {
		errRes = TokenNotExist
	} else {
		token = kv[1]

		if claims, err := ParseToken(token);err != nil {
			errRes = TokenValidateFailed
		} else {
			if time.Now().Unix() > claims.ExpiresAt {
				errRes = TokenExpired
			}
		}
	}

	fmt.Println("error:", errRes)
	if errRes != nil {
		responseWriter := response.Default(c)
		responseWriter.Error(401, errRes.Error())

		c.Abort()
		return
	}

	c.Next()
}

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
	"github.com/ebar-go/ego/http/constant"
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
	Verification int `json:"verification"`
}

var TokenNotExist = errors.New("token not exist")
var TokenValidateFailed = errors.New("token validate failed")
var TokenExpired = errors.New("token expired")

// SetJwtSecret 设置jwt的秘钥
func SetJwtSecret(secret []byte) {
	jwtSecret = secret
}

// GetEncodeToken 获取加密的token
func GetEncodeToken(iss, secret string, expireTime int) (string, error) {
	if expireTime == 0 {
		expireTime = constant.JwtExpiredTime
	}
	now := time.Now().Unix()
	exp := now + int64(expireTime)
	claim := jwt.MapClaims{
		"iss":       iss,
		"iat":      now,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claim)
	tokenStr,err  := token.SignedString([]byte(secret))
	return tokenStr, err
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

// GetCurrentUser 获取解析jwt后的用户信息
func GetCurrentUser(ctx *gin.Context) interface{} {
	user, exist := ctx.Get(constant.JwtUserKey)
	if !exist {
		return nil
	}

	return user
}


// JWT gin的jwt中间件
func JWT(c *gin.Context) {
	var errRes error

	// 获取token
	token := c.GetHeader(constant.JwtTokenHeader)
	fmt.Println(token)

	kv := strings.Split(token, " ")
	if len(kv) != 2 || kv[0] != constant.JwtTokenMethod {
		errRes = TokenNotExist
	} else {
		token = kv[1]

		claims, err := ParseToken(token)
		if err != nil {
			errRes = TokenValidateFailed
			fmt.Println(err)
		} else {
			if time.Now().Unix() > claims.ExpiresAt {
				errRes = TokenExpired
			}else {
				// 存储用户信息
				c.Set(constant.JwtUserKey, &claims.User)
			}
		}
	}

	if errRes != nil {
		response.Error(c, constant.StatusUnauthorized, errRes.Error())

		c.Abort()
		return
	}


	c.Next()
}

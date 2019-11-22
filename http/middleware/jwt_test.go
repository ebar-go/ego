package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwtAuth_CreateToken(t *testing.T) {
	auth := JwtAuth{
		SigningKey: []byte("test"),
	}

	type MyClaims struct {
		User struct {
			UID  int    `json:"uid"`
			ACID string `json:"acId"`
			Name string `json:"name"`
			Verification int `json:"verification"`
		}
		jwt.StandardClaims
	}

	token , err := auth.CreateToken(func() jwt.Claims {
		claims := MyClaims{}
		claims.User.UID = 1
		claims.Subject = "test"
		claims.ExpiresAt = time.Now().Unix() + 10

		return claims
	})

	assert.Nil(t, err)
	fmt.Println(token)

	claims, err := auth.parseToken(token)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println(claims)


}

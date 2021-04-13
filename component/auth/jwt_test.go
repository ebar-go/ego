package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customClaims struct {
	jwt.StandardClaims
	Id int
}

func TestJwtImpl_GenerateToken(t *testing.T) {
	c := new(customClaims)
	c.Id = 1
	token, err := NewJwt([]byte("123456")).GenerateToken(c)
	assert.Nil(t, err)
	fmt.Println(token)
}

func TestJwtImpl_ParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MX0.O9eeIgDqrqA_otMywvqkiyUOd0m1UzXXig0Xs6ikQ9o"
	c, err := NewJwt([]byte("123456")).ParseToken(token)
	assert.Nil(t, err)
	fmt.Println(c)
}

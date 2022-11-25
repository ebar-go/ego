package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

// Instance provide the jwt algorithm component
type Instance struct{}

// CreateToken 生成token
func (impl Instance) GenerateToken(sign []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sign)
}

// ParseToken return jwt.MapClaims and error
func (impl Instance) ParseToken(sign []byte, token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return sign, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, ErrTokenInvalid
	}

	return tokenClaims.Claims, nil
}

func New() *Instance {
	return &Instance{}
}

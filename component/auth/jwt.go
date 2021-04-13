package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type Jwt interface {
	// 生成token
	GenerateToken(claims jwt.Claims) (string, error)
	// 解析token
	ParseToken(token string) (jwt.Claims, error)
}

// jwtImpl jwt
type jwtImpl struct {
	key []byte
}

// CreateToken 生成token
func (impl jwtImpl) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(impl.key)
}

// ParseToken return jwt.MapClaims and error
func (impl jwtImpl) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return impl.key, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, errors.New("token is invalid")
	}

	return tokenClaims.Claims, nil
}

// NewJwt return jwtImpl instance
func NewJwt(key []byte) *jwtImpl {
	return &jwtImpl{key: key}
}

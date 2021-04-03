package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type Jwt interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ParseToken(token string) (jwt.Claims, error)
	ParseTokenWithClaims(token string, claims jwt.Claims) error
}

// JwtAuth jwt
type JwtAuth struct {
	key   []byte
}

var (
	TokenValidateFailed = errors.New("token validate failed")
)

// CreateToken 生成token
func (jwtAuth JwtAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAuth.key)
}

// ParseToken parse token
func (jwtAuth JwtAuth) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.key, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, TokenValidateFailed
	}

	return tokenClaims.Claims, nil
}

// ParseToken parse token
func (jwtAuth JwtAuth) ParseTokenWithClaims(token string, claims jwt.Claims) error {
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.key, nil
	})

	if err != nil {
		return err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return TokenValidateFailed
	}

	claims = tokenClaims.Claims
	return nil
}

// NewJwt return JwtAuth instance
func NewJwt(key []byte) Jwt {
	return &JwtAuth{key: key}
}
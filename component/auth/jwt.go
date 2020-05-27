package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ebar-go/ego/utils/number"
	"time"
)

// Jwt json web token
type Jwt interface {
	// Parse token
	ParseToken(token string) (jwt.Claims, error)

	// Create token
	CreateToken(claimsCreator func() jwt.Claims) (string, error)

	GenerateToken(tokenExpireTime int, iss string) (string, error)
}

// JwtAuth jwt
type JwtAuth struct {
	SignKey []byte
}

// New return JwtAuth instance
func New(signKey []byte) Jwt {
	return &JwtAuth{SignKey: signKey}
}

var (
	TokenValidateFailed = errors.New("token validate failed")
)

// CreateToken create token
func (jwtAuth JwtAuth) CreateToken(claimsCreator func() jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsCreator())
	return token.SignedString(jwtAuth.SignKey)
}

// ParseToken parse token
func (jwtAuth JwtAuth) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.SignKey, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, TokenValidateFailed
	}

	return tokenClaims.Claims, nil
}

// GenerateToken
func (jwtAuth JwtAuth) GenerateToken(tokenExpireTime int, iss string) (string, error) {
	tokenExpireTime = number.DefaultInt(tokenExpireTime, 3000)
	now := time.Now().Unix()
	exp := now + int64(tokenExpireTime)
	claim := jwt.MapClaims{
		"iss": iss,
		"iat": now,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(jwtAuth.SignKey)
	return tokenStr, err
}

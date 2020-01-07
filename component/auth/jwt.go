package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// Jwt json web token
type Jwt interface {
	// Parse token
	ParseToken(token string) (jwt.Claims, error)

	// Create token
	CreateToken(claimsCreator func() jwt.Claims) (string, error)
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

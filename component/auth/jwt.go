package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// JwtAuth jwt
type JwtAuth struct {
	signKey   []byte
	ClaimsKey string
}

const (
	defaultClaimsKey = "jwt_claims"
)

// New return JwtAuth instance
func New(signKey []byte) *JwtAuth {
	return &JwtAuth{signKey: signKey, ClaimsKey: defaultClaimsKey}
}

var (
	TokenValidateFailed = errors.New("token validate failed")
)

// CreateToken 生成token
func (jwtAuth JwtAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAuth.signKey)
}

// ParseToken parse token
func (jwtAuth JwtAuth) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.signKey, nil
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
		return jwtAuth.signKey, nil
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

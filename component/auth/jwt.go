package auth

import (
	"errors"
)

// Jwt json web token
type Jwt interface {
	// 解析token
	ParseToken(token string) (jwt.Claims, error)

	// 创建token
	CreateToken(claimsCreator func() jwt.Claims) (string, error)
}

func NewJwt(signKey []byte) Jwt {
	return &JwtAuth{SignKey: signKey}
}

// JwtAuth jwt
type JwtAuth struct {
	SignKey []byte
}

var (
	TokenValidateFailed = errors.New("token validate failed")
)

// CreateToken 生成一个token
func (jwtAuth JwtAuth) CreateToken(claimsCreator func() jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsCreator())
	return token.SignedString(jwtAuth.SignKey)
}

// ParseToken 解析Token
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

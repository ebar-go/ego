package component

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

// JWT provide the jwt algorithm component
type JWT struct {
	Named
	key []byte
}

// CreateToken 生成token
func (impl JWT) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(impl.key)
}

// ParseToken return jwt.MapClaims and error
func (impl JWT) ParseToken(token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return impl.key, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Claims == nil || !tokenClaims.Valid {
		return nil, ErrTokenInvalid
	}

	return tokenClaims.Claims, nil
}

func NewJWT() *JWT {
	return &JWT{Named: "jwt"}
}
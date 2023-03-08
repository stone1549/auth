package service

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

// TokenFactory provides methods for creating authentication tokens.
type TokenFactory interface {
	// NewToken returns a new token string with the given claims
	NewToken(claims Claims) (string, error)
}

type Claims struct {
	// Subject (globally unique user id) of token
	Sub string

	// Subjects email address
	Email string

	// Subjects username
	Username string

	// Not valid before
	Nbf int64

	// Expire at
	Exp int64

	// Issued at
	Iat int64
}

func NewClaims(id, email, username string) Claims {
	now := time.Now().Unix()
	exp := time.Now().Add(time.Hour).Unix()
	return Claims{id, email, username, now, exp, now}
}

type jwtFactory struct {
	SigningMethod   jwt.SigningMethod
	SecretSharedKey []byte
	RsaPrivateKey   *rsa.PrivateKey
	RsaPublicKey    *rsa.PublicKey
}

// NewToken returns a new token string with the given claims
func (jwtf *jwtFactory) NewToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwtf.SigningMethod, jwt.MapClaims{
		"sub":      claims.Sub,
		"email":    claims.Email,
		"username": claims.Username,
		"nbf":      claims.Nbf,
		"exp":      claims.Exp,
		"iat":      claims.Iat,
	})

	if jwtf.SigningMethod == jwt.SigningMethodRS512 {
		return token.SignedString(jwtf.RsaPrivateKey)
	} else if jwtf.SigningMethod == jwt.SigningMethodHS512 {
		return token.SignedString([]byte(jwtf.SecretSharedKey))
	} else {
		return "", errors.New("unsupported JWT configuration")
	}
}

// NewTokenFactory constructs a token factory using the given configuration.
func NewTokenFactory(config Configuration) (TokenFactory, error) {
	if config.GetTokenSecretKey() != "" {
		return &jwtFactory{
			jwt.SigningMethodHS512,
			[]byte(config.GetTokenSecretKey()),
			nil,
			nil}, nil
	} else if config.GetTokenPublicKey() != nil && config.GetTokenPrivateKey() != nil {
		return &jwtFactory{
			jwt.SigningMethodRS512,
			nil,
			config.GetTokenPrivateKey(),
			config.GetTokenPublicKey()}, nil
	} else {
		return nil, errors.New("invalid token signing configuration")
	}
}

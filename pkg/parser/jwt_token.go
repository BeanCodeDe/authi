package parser

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/golang-jwt/jwt"
)

const (
	// Environment variable to point to path of the public key file
	EnvPublicKeyPath = "PUBLIC_KEY_PATH"
)

var (
	// Token could not be recognized in authorizationString
	ErrTokenNotFound = errors.New("token not found")
	// Claim doesn't match content from token
	ErrClaimCouldNotBeParsed = errors.New("claim could not be parsed")
	// Token is not valid. Maybe it is expired or doesn't match the public key
	ErrTokenNotValid = errors.New("token is not valid")
	// Public key from environment variable PUBLIC_KEY_PATH couldn't be read
	ErrWhileReadingKey = errors.New("error while reading public key")
	// Public key seams to not have the valid format
	ErrWhileParsingKey = errors.New("error while parsing public key")
)

// Struct to parse and validate jwt tokens
type JWTParser struct {
	verifyKey *rsa.PublicKey
}

// Constructor to create jwt parser. public key path have to be set under environment variable PUBLIC_KEY_PATH
func NewJWTParser() (Parser, error) {
	publicKeyPath := util.GetEnvWithFallback(EnvPublicKeyPath, "/token/jwtRS256.key.pub")

	verifyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWhileReadingKey, err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWhileParsingKey, err)
	}
	return &JWTParser{verifyKey: verifyKey}, nil
}

// Checks if the bearer token is valid and returns its content as a claim
func (parser *JWTParser) ParseToken(authorizationString string) (*adapter.Claims, error) {
	splitToken := strings.Split(authorizationString, "Bearer ")
	if len(splitToken) != 2 {
		return nil, ErrTokenNotFound
	}
	tokenString := splitToken[1]

	claims := &adapter.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return parser.verifyKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrClaimCouldNotBeParsed, err)
	}

	if tkn == nil || !tkn.Valid {
		return nil, fmt.Errorf("%w: %w", ErrTokenNotValid, err)
	}

	return claims, nil
}

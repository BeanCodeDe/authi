package authadapter

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	AuthorizationHeaderName string = "Authorization"
	RefreshTokenHeaderName  string = "refresh_token"
	ClaimName               string = "claim"
)

type (
	Auth interface {
		ParseToken(authorizationString string) (*Claims, error)
		CreateJWTToken(token string) (string, error)
	}
	Claims struct {
		UserId uuid.UUID `json:"user_id"`
		jwt.StandardClaims
	}
	AuthAdapter struct {
		verifyKey     *rsa.PublicKey
		publicKeyPath string
		authUrl       string
	}
)

func NewAuthAdapter() (*AuthAdapter, error) {
	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	authUrl := os.Getenv("AUTH_URL")
	verifyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading public Key: %v", err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, fmt.Errorf("error while parsing public Key: %v", err)
	}

	return &AuthAdapter{verifyKey: verifyKey, publicKeyPath: publicKeyPath, authUrl: authUrl}, nil
}

func (authAdapter *AuthAdapter) ParseToken(authorizationString string) (*Claims, error) {
	splitToken := strings.Split(authorizationString, "Bearer ")
	if len(splitToken) != 2 {
		return nil, fmt.Errorf("token not found")
	}
	tokenString := splitToken[1]

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return authAdapter.verifyKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("claim could not be parsed: %v", err)
	}

	if tkn == nil || !tkn.Valid {
		return nil, fmt.Errorf("token is not valid: %v", err)
	}

	return claims, nil
}

func (authAdapter *AuthAdapter) CreateJWTToken(token string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, authAdapter.authUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set(AuthorizationHeaderName, token)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status is no ok but %v", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("body with token couldn't be read: %v", err)
	}

	stringToken := string(bodyBytes)

	return stringToken, nil
}

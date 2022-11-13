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

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

var (
	verifyKey      *rsa.PublicKey
	pubblicKeyPath string
	authUrl        string
)

func Init() error {
	pubblicKeyPath = os.Getenv("PUBLIC_KEY_PATH")
	authUrl = os.Getenv("AUTH_URL")
	verifyBytes, err := os.ReadFile(pubblicKeyPath)
	if err != nil {
		return fmt.Errorf("error while reading public Key: %v", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return fmt.Errorf("error while parsing public Key: %v", err)
	}

	return nil
}

func ParseToken(authorizationString string) (*Claims, error) {
	splitToken := strings.Split(authorizationString, "Bearer ")
	tokenString := splitToken[1]

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("claim could not be parsed: %v", err)
	}

	if tkn == nil || !tkn.Valid {
		return nil, fmt.Errorf("token is not valid: %v", err)
	}

	return claims, nil
}

func CreateJWTToken(token string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
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
		return "", fmt.Errorf("body with token coulnd be read: %v", err)
	}

	stringToken := string(bodyBytes)

	return stringToken, nil
}

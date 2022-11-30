package util

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/golang-jwt/jwt"
	"gopkg.in/go-playground/assert.v1"
)

const (
	loginPath   = "/login"
	refreshPath = "/refresh"
)

const (
	PublicKeyFile      = "./data/token/public/jwtRS256.key.pub"
	PrivatKeyFile      = "./data/token/privat/jwtRS256.key"
	WrongPublicKeyFile = "./data/token/public/jwtRS256_wrong.key.pub"
	WrongPrivatKeyFile = "./data/token/privat/jwtRS256_wrong.key"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func sendLoginRequest(userId string, authenticate *Authenticate) *http.Response {
	userJson, err := json.Marshal(authenticate)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url+userPath+"/"+userId+loginPath, bytes.NewBuffer(userJson))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", adapter.ContentTyp)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func sendRefreshTokenRequest(userId string, token string, refreshToken string) *http.Response {
	req, err := http.NewRequest(http.MethodPatch, url+userPath+"/"+userId+refreshPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(adapter.AuthorizationHeaderName, "Bearer "+token)
	req.Header.Set(adapter.RefreshTokenHeaderName, refreshToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func Login(userId string, authenticate *Authenticate) (*TokenResponseDTO, int) {
	response := sendLoginRequest(userId, authenticate)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode
	}
	token := new(TokenResponseDTO)
	if err := json.NewDecoder(response.Body).Decode(token); err != nil {
		panic(err)
	}
	return token, response.StatusCode
}

func RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, int) {
	response := sendRefreshTokenRequest(userId, token, refreshToken)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode
	}
	tokenResponse := new(TokenResponseDTO)
	if err := json.NewDecoder(response.Body).Decode(tokenResponse); err != nil {
		panic(err)
	}
	return tokenResponse, response.StatusCode
}

func ObtainToken(t *testing.T) (*TokenResponseDTO, string, *Authenticate) {
	userId := CreateUserForFurtherTesting(t)
	authenticate := &Authenticate{Password: DefaultPassword}
	token, status := Login(userId, authenticate)
	assert.Equal(t, status, http.StatusOK)
	return token, userId, authenticate
}

func CreateCustomJWTToken(userId string, expirationTime int64, signKey *rsa.PrivateKey) string {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		panic(err)
	}

	return signedToken
}

func LoadPrivatKeyFile(fileName string) *rsa.PrivateKey {
	verifyBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	verifyKey, err := jwt.ParseRSAPrivateKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}

	return verifyKey
}

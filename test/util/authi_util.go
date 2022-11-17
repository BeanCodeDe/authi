package util

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/golang-jwt/jwt"
	"gopkg.in/go-playground/assert.v1"
)

const (
	authPath    = "/auth"
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

func sendLoginRequest(user *UserDTO) *http.Response {
	userJson, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodGet, url+authPath+loginPath, bytes.NewBuffer(userJson))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentTyp)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func sendRefreshTokenRequest(token string, refreshToken string) *http.Response {
	req, err := http.NewRequest(http.MethodPatch, url+authPath+refreshPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(authadapter.AuthorizationHeaderName, "Bearer "+token)
	req.Header.Set(authadapter.RefreshTokenHeaderName, refreshToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func Login(loginUser *UserDTO) (*TokenResponseDTO, int) {
	response := sendLoginRequest(loginUser)
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

func RefreshToken(token string, refreshToken string) (*TokenResponseDTO, int) {
	response := sendRefreshTokenRequest(token, refreshToken)
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

func OptainToken(t *testing.T) (*TokenResponseDTO, *UserDTO) {
	userId := CreateUserForFurtherTesting(t)
	user := &UserDTO{ID: userId, Password: DefaultPassword}
	token, status := Login(user)
	assert.Equal(t, status, http.StatusOK)
	return token, user
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

package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"gopkg.in/go-playground/assert.v1"
)

const (
	authPath    = "/auth"
	loginPath   = "/login"
	refreshPath = "/refresh"
)

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

func OptainToken(t *testing.T) *TokenResponseDTO {
	userId := CreateUserForFurtherTesting(t)
	user := &UserDTO{ID: userId, Password: DefaultPassword}
	token, status := Login(user)
	assert.Equal(t, status, http.StatusOK)
	return token
}

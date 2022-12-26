package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	// Name of header where jwt token is stored
	AuthorizationHeaderName = "Authorization"
	// Name of header where refresh token is stored
	RefreshTokenHeaderName = "refresh_token"
	// Key to set claim in context
	ClaimName = "claim"
	// Environment variable to point to url for the auth service
	EnvAuthUrl = "AUTHI_URL"

	// Url of root path for authi
	AuthiRootPath = "/user"
	//Path to login api
	AuthiLoginPath = "/login"
	//Path to refresh api
	AuthiRefreshPath = "/refresh"

	//Content type for AuthenticateDTO
	ContentTyp    = "application/json; charset=utf-8"
	correlationId = "X-Correlation-ID"
)

var (
	//Error that indicates, that the returned http status is not 200
	errStatusNotOk = errors.New("status is no ok")
	//Error that indicates, that the returned body can not be parsed into TokenResponseDTO
	errReadResponse = errors.New("body with token couldn't be read")
)

// Implementation of auth adapter to login and refresh token of user
type AuthiAdapter struct {
	authiRefreshUrl string
	authiLoginUrl   string
	correlationId   string
}

// Initialize auth adapter with public key and url to authi service.
// Therefor the environment variable AUTH_URL have to be set
func NewAuthiAdapter(correlationId string) AuthAdapter {
	authUrl := os.Getenv(EnvAuthUrl)
	authiRefreshUrl := authUrl + AuthiRootPath + "/%s" + AuthiRefreshPath
	authiLoginUrl := authUrl + AuthiRootPath + "/%s" + AuthiLoginPath
	return &AuthiAdapter{authiRefreshUrl: authiRefreshUrl, authiLoginUrl: authiLoginUrl, correlationId: correlationId}
}

// Get new token with refresh token from authi service
func (authAdapter *AuthiAdapter) RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf(authAdapter.authiRefreshUrl, userId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(AuthorizationHeaderName, "Bearer "+token)
	req.Header.Set(RefreshTokenHeaderName, refreshToken)
	req.Header.Set(correlationId, uuid.NewString())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return readTokenResponse(resp)
}

// Login to get token
func (authAdapter *AuthiAdapter) GetToken(userId string, password string) (*TokenResponseDTO, error) {
	authenticateJson, err := json.Marshal(&AuthenticateDTO{Password: password})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(authAdapter.authiLoginUrl, userId), bytes.NewBuffer(authenticateJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ContentTyp)
	req.Header.Set(correlationId, uuid.NewString())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return readTokenResponse(resp)
}

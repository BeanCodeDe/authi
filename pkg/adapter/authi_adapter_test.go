package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_Successfully(t *testing.T) {
	userId := uuid.New()
	token := "someToken"
	refreshToken := "someRefreshToken"
	tokenResponse := &TokenResponseDTO{}
	// Setup
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Contains(t, req.URL.Path, userId.String())
		assert.Equal(t, "Bearer "+token, req.Header.Get(AuthorizationHeaderName))
		assert.Equal(t, refreshToken, req.Header.Get(RefreshTokenHeaderName))

		tokenResponseJSON, err := json.Marshal(&tokenResponse)
		assert.Nil(t, err)

		res.Write(bytes.NewBuffer(tokenResponseJSON).Bytes())
	}))
	defer func() { testServer.Close() }()
	authAdapter := getAuthiAdapter(testServer.URL)
	// Exec
	result, err := authAdapter.RefreshToken(userId.String(), token, refreshToken)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, tokenResponse, result)
}

func TestRefreshToken_ErrorWhileCreatingRequest(t *testing.T) {
	userId := uuid.New()
	token := "someToken"
	refreshToken := "someRefreshToken"
	// Setup
	authAdapter := &AuthiAdapter{}

	// Exec
	result, err := authAdapter.RefreshToken(userId.String(), token, refreshToken)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestRefreshToken_ErrorWhileExecutingRequest(t *testing.T) {
	userId := uuid.New()
	token := "someToken"
	refreshToken := "someRefreshToken"
	tokenResponse := &TokenResponseDTO{}
	// Setup
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Contains(t, req.URL.Path, userId.String())
		assert.Equal(t, "Bearer "+token, req.Header.Get(AuthorizationHeaderName))
		assert.Equal(t, refreshToken, req.Header.Get(RefreshTokenHeaderName))

		tokenResponseJSON, err := json.Marshal(&tokenResponse)
		assert.Nil(t, err)

		res.Write(bytes.NewBuffer(tokenResponseJSON).Bytes())
	}))
	defer func() { testServer.Close() }()
	authAdapter := getAuthiAdapter(testServer.URL + "somethingWrong")
	// Exec
	result, err := authAdapter.RefreshToken(userId.String(), token, refreshToken)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestGetToken_Successfully(t *testing.T) {
	userId := uuid.New()
	password := "password"
	tokenResponse := &TokenResponseDTO{}
	// Setup
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Contains(t, req.URL.Path, userId.String())
		assert.Equal(t, ContentTyp, req.Header.Get("Content-Type"))

		tokenResponseJSON, err := json.Marshal(&tokenResponse)
		assert.Nil(t, err)

		res.Write(bytes.NewBuffer(tokenResponseJSON).Bytes())
	}))
	defer func() { testServer.Close() }()
	authAdapter := getAuthiAdapter(testServer.URL)
	// Exec
	result, err := authAdapter.GetToken(userId.String(), password)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, tokenResponse, result)
}

func TestGetToken_ErrorWhileParsingPassword(t *testing.T) {
	userId := uuid.New()
	password := ""

	// Setup
	authAdapter := &AuthiAdapter{}

	// Exec
	result, err := authAdapter.GetToken(userId.String(), password)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestGetToken_ErrorWhileCreatingRequest(t *testing.T) {
	userId := uuid.New()
	password := "someToken"

	// Setup
	authAdapter := &AuthiAdapter{}

	// Exec
	result, err := authAdapter.GetToken(userId.String(), password)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestGetToken_ErrorWhileExecutingRequest(t *testing.T) {
	userId := uuid.New()
	password := "password"
	tokenResponse := &TokenResponseDTO{}
	// Setup
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Contains(t, req.URL.Path, userId.String())
		assert.Equal(t, ContentTyp, req.Header.Get("Content-Type"))

		tokenResponseJSON, err := json.Marshal(&tokenResponse)
		assert.Nil(t, err)

		res.Write(bytes.NewBuffer(tokenResponseJSON).Bytes())
	}))
	defer func() { testServer.Close() }()
	authAdapter := getAuthiAdapter(testServer.URL + "somethingWrong")
	// Exec
	result, err := authAdapter.GetToken(userId.String(), password)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func getAuthiAdapter(url string) AuthAdapter {
	authAdapter := &AuthiAdapter{
		authiRefreshUrl: url + AuthiRootPath + "/%s" + AuthiRefreshPath,
		authiLoginUrl:   url + AuthiRootPath + "/%s" + AuthiLoginPath,
	}
	return authAdapter
}

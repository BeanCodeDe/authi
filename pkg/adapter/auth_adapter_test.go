package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	RefreshTokenRecord struct {
		userId       string
		token        string
		refreshToken string
	}

	GetTokenRecord struct {
		userId   string
		password string
	}

	RefreshTokenResponse struct {
		tokenResponseDTO *TokenResponseDTO
		err              error
	}

	GetTokenResponse struct {
		tokenResponseDTO *TokenResponseDTO
		err              error
	}

	AdapterMock struct {
		refreshTokenRecordArray   []*RefreshTokenRecord
		refreshTokenResponseArray []*RefreshTokenResponse

		getTokenRecordArray   []*GetTokenRecord
		getTokenResponseArray []*GetTokenResponse
	}
)

func (mock *AdapterMock) RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, error) {
	refreshTokenRecord := &RefreshTokenRecord{userId: userId, token: token, refreshToken: refreshToken}
	mock.refreshTokenRecordArray = append(mock.refreshTokenRecordArray, refreshTokenRecord)

	response := mock.refreshTokenResponseArray[len(mock.refreshTokenResponseArray)-1]
	return response.tokenResponseDTO, response.err
}
func (mock *AdapterMock) GetToken(userId string, password string) (*TokenResponseDTO, error) {
	getTokenRecord := &GetTokenRecord{userId: userId, password: password}
	mock.getTokenRecordArray = append(mock.getTokenRecordArray, getTokenRecord)

	response := mock.getTokenResponseArray[len(mock.getTokenResponseArray)-1]
	return response.tokenResponseDTO, response.err
}

func TestReadTokenResponse_Successfully(t *testing.T) {
	// Setup
	tokenResponse := &TokenResponseDTO{}
	tokenResponseJSON, err := json.Marshal(&tokenResponse)
	assert.Nil(t, err)
	byteResponse := bytes.NewBuffer(tokenResponseJSON)

	respRecorder := &httptest.ResponseRecorder{Code: http.StatusOK, Body: byteResponse}
	resp := respRecorder.Result()

	// Exec
	result, err := readTokenResponse(resp)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, tokenResponse, result)
}

func TestReadTokenResponse_WrongStatusCode(t *testing.T) {
	// Setup
	respRecorder := &httptest.ResponseRecorder{Code: http.StatusConflict}
	resp := respRecorder.Result()

	// Exec
	tokenResponse, err := readTokenResponse(resp)

	// Assertions
	assert.Nil(t, tokenResponse)
	assert.ErrorIs(t, err, errStatusNotOk)
}

func TestReadTokenResponse_WrongBodyObject(t *testing.T) {
	// Setup
	respRecorder := &httptest.ResponseRecorder{Code: http.StatusOK, Body: bytes.NewBufferString("some data")}
	resp := respRecorder.Result()

	// Exec
	tokenResponse, err := readTokenResponse(resp)

	// Assertions
	assert.Nil(t, tokenResponse)
	assert.ErrorIs(t, err, errReadResponse)
}

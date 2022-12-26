package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

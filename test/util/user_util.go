package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

const (
	createUserJson  = `{"password":"%s"}`
	DefaultPassword = "SomeDefaultPassowrd"
)

func sendCreateUserIdRequest() *http.Response {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, Url+adapter.AuthiRootPath, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set(CorrelationId, uuid.NewString())
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	return resp
}

func sendCreateUserRequest(id string, userCreateJson string) *http.Response {
	jsonReq := []byte(userCreateJson)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, Url+adapter.AuthiRootPath+"/"+id, bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", adapter.ContentTyp)
	req.Header.Set(CorrelationId, uuid.NewString())
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	return resp
}

func sendRefreshPasswordRequest(userId string, authenticate *Authenticate, token string) *http.Response {
	userJson, err := json.Marshal(authenticate)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPatch, Url+adapter.AuthiRootPath+"/"+userId, bytes.NewBuffer(userJson))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", adapter.ContentTyp)
	req.Header.Set(adapter.AuthorizationHeaderName, "Bearer "+token)
	req.Header.Set(CorrelationId, uuid.NewString())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func sendDeleteUserRequest(userId string, token string) *http.Response {
	req, err := http.NewRequest(http.MethodDelete, Url+adapter.AuthiRootPath+"/"+userId, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(adapter.AuthorizationHeaderName, "Bearer "+token)
	req.Header.Set(CorrelationId, uuid.NewString())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func CreateUserId() (string, int) {
	response := sendCreateUserIdRequest()
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return "", response.StatusCode
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not parse response Body: %v", err)
		return "", -1
	}

	uuidOfUser := string(bodyBytes)

	return uuidOfUser, response.StatusCode
}

func CreateUser(id string, password string) int {
	response := sendCreateUserRequest(id, fmt.Sprintf(createUserJson, password))
	defer response.Body.Close()
	return response.StatusCode
}

func CreateUserForFurtherTesting(t *testing.T) string {
	userId, status := CreateUserId()
	assert.Equal(t, status, http.StatusCreated)
	_, err := uuid.Parse(userId)
	assert.Equal(t, err, nil)

	status = CreateUser(userId, DefaultPassword)
	assert.Equal(t, status, http.StatusCreated)
	return userId
}

func UpdatePassword(userId string, password string, tokenString string) int {
	response := sendRefreshPasswordRequest(userId, &Authenticate{Password: password}, tokenString)
	defer response.Body.Close()
	return response.StatusCode
}

func DeleteUser(userId string, tokenString string) int {
	response := sendDeleteUserRequest(userId, tokenString)
	defer response.Body.Close()
	return response.StatusCode
}

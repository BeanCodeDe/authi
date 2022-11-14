package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

const (
	userPath        = "/user"
	createUserJson  = `{"password":"%s"}`
	DefaultPassword = "SomeDefaultPassowrd"
)

func sendCreateUserIdRequest() *http.Response {
	resp, err := http.Post(url+userPath, contentTyp, nil)
	if err != nil {
		panic(err)
	}

	return resp
}

func sendCreateUserRequest(id string, userCreateJson string) *http.Response {
	jsonReq := []byte(userCreateJson)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url+userPath+"/"+id, bytes.NewBuffer(jsonReq))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", contentTyp)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	return resp
}

func CreateUserId() (string, int) {
	response := sendCreateUserIdRequest()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
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
	assert.Equal(t, status, http.StatusOK)
	_, err := uuid.Parse(userId)
	assert.Equal(t, err, nil)

	status = CreateUser(userId, DefaultPassword)
	assert.Equal(t, status, http.StatusCreated)
	return userId
}

package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	authPath  = "/auth"
	loginPath = "/login"
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

func Login(loginUser *UserDTO) (*TokenResponseDTO, int) {
	response := sendLoginRequest(loginUser)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode
	}
	token := new(TokenResponseDTO)
	json.NewDecoder(response.Body).Decode(token)
	return token, response.StatusCode
}

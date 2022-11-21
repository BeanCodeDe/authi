package core

import (
	"math/rand"

	"github.com/google/uuid"
)

type (
	TokenResponseDTO struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}

	AuthenticateDTO struct {
		Password string `json:"password" validate:"required"`
	}

	Facade interface {
		CreateUser(userId uuid.UUID, authenticate *AuthenticateDTO) error
		LoginUser(userId uuid.UUID, authenticate *AuthenticateDTO) (*TokenResponseDTO, error)
		RefreshToken(userId uuid.UUID, refreshToken string) (*TokenResponseDTO, error)
	}
)

const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func randomString() string {
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	return string(bytes)
}

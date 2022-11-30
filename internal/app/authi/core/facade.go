package core

import (
	"math/rand"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
)

type (
	Facade interface {
		CreateUser(userId uuid.UUID, authenticate *adapter.AuthenticateDTO) error
		LoginUser(userId uuid.UUID, authenticate *adapter.AuthenticateDTO) (*adapter.TokenResponseDTO, error)
		RefreshToken(userId uuid.UUID, refreshToken string) (*adapter.TokenResponseDTO, error)
		UpdatePassword(userId uuid.UUID, authenticate *adapter.AuthenticateDTO) error
		DeleteUser(userId uuid.UUID) error
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

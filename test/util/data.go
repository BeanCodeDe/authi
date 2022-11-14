package util

const (
	url = "http://localhost:1203"

	contentTyp = "application/json; charset=utf-8"
)

type (
	UserDTO struct {
		ID       string `json:"id" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	TokenResponseDTO struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}
)

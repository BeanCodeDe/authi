package util

const (
	Url = "http://localhost:1203"
)

type (
	Authenticate struct {
		Password string `json:"password" validate:"required"`
	}

	TokenResponseDTO struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}
)

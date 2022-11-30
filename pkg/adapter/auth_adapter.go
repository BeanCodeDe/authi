// Package to request and verify jwt token
package adapter

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	//Interface for authentication
	AuthAdapter interface {
		RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, error)
		GetToken(userId string, password string) (*TokenResponseDTO, error)
	}
	//Claim with data from the token
	Claims struct {
		UserId uuid.UUID `json:"user_id"`
		jwt.StandardClaims
	}
	//Response with token data
	TokenResponseDTO struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}
	//Request object for authentication
	AuthenticateDTO struct {
		Password string `json:"password" validate:"required"`
	}
)

package api

import (
	"net/http"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const AuthRootPath = "/auth"

type tokenResponseDTO struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

func InitAuthInterface(group *echo.Group) {
	group.PATCH("/refresh", refreshToken)
	group.GET("/login", login)
}

func refreshToken(context echo.Context) error {
	log.Debugf("Refresh token")

	refreshToken := context.Request().Header.Get(authadapter.RefreshTokenHeaderName)
	claims, ok := context.Get(authadapter.ClaimName).(authadapter.Claims)
	if !ok {
		log.Errorf("Got data of wrong type: %v", context.Get(authadapter.ClaimName))
		return echo.ErrUnauthorized
	}

	token, err := core.CreateJWTTokenFromRefreshToken(claims.UserId, refreshToken)
	if err != nil {
		log.Errorf("Something went wrong while creating Token: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("Refresh token for user %s updated", claims.UserId)
	return context.JSON(http.StatusOK, token)
}

func login(context echo.Context) error {
	log.Debugf("Login some user")
	userCore, err := bind(context, new(userLoginDTO))
	if err != nil {
		log.Warnf("Error while binding user: %v", err)
		return echo.ErrBadRequest
	}

	tokenCore, err := userCore.Login()
	if err != nil {
		log.Warnf("Error while logging in user %v: %v", userCore, err)
		return echo.ErrUnauthorized
	}

	log.Debugf("Logged in user %s", userCore.ID)
	return context.JSON(http.StatusOK, mapToTokenResponseDTO(tokenCore))
}

func mapToTokenResponseDTO(tokenCore *core.TokenCore) *tokenResponseDTO {
	return &tokenResponseDTO{AccessToken: tokenCore.AccessToken, ExpiresIn: tokenCore.ExpiresIn, RefreshToken: tokenCore.RefreshToken, RefreshExpiresIn: tokenCore.RefreshExpiresIn}
}

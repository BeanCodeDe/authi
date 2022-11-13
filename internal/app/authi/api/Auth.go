package api

import (
	"net/http"

	"github.com/BeanCodeDe/SpaceLight-Auth/internal/auth"
	"github.com/BeanCodeDe/SpaceLight-AuthMiddleware/authAdapter"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const AuthRootPath = "/auth"

func InitAuthInterface(group *echo.Group) {
	group.GET("/refresh", refreshToken, auth.AuthMiddleware)
}

func refreshToken(context echo.Context) error {
	log.Debugf("Refresh token")
	claims, ok := context.Get(authAdapter.ClaimName).(authAdapter.Claims)
	if !ok {
		log.Errorf("Got data of wrong type: %v", context.Get(authAdapter.ClaimName))
		return echo.ErrUnauthorized
	}

	token, err := auth.CreateJWTToken(claims.UserId, claims.Roles)
	if err != nil {
		log.Errorf("Something went wrong while creating Token: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("Logged in user %s", claims.UserId)
	return context.String(http.StatusOK, token)
}

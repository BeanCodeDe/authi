package api

import (
	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type (
	Api interface {
		CreateUserId(context echo.Context) error
		CreateUser(context echo.Context) error
		RefreshToken(context echo.Context) error
		LoginUser(context echo.Context) error
	}
)

const (
	userRootPath    = "/user"
	userIdParam     = "userId"
	userLoginPath   = "/login"
	userRefreshPath = "/refresh"
)

func bindAuthenticate(context echo.Context) (uuid.UUID, *core.AuthenticateDTO, error) {
	log.Debugf("Bind context to auth %v", context)
	authenticate := new(core.AuthenticateDTO)
	if err := context.Bind(authenticate); err != nil {
		log.Warnf("Could not bind auth, %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}
	log.Debugf("Auth bind %v", authenticate)
	if err := context.Validate(authenticate); err != nil {
		log.Warnf("Could not validate auth, %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}

	return userId, authenticate, nil
}

func checkUserId(context echo.Context, userId uuid.UUID) error {
	claims, ok := context.Get(authadapter.ClaimName).(authadapter.Claims)

	if !ok {
		log.Warnf("Could not map Claims")
		return echo.ErrUnauthorized
	}

	if userId != claims.UserId {
		log.Warnf("User %v is not allowed to get token for user %v", userId, claims.UserId)
		return echo.ErrUnauthorized
	}

	return nil
}

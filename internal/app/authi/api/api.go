package api

import (
	"github.com/BeanCodeDe/authi/pkg/adapter"
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

func bindAuthenticate(context echo.Context) (uuid.UUID, *adapter.AuthenticateDTO, error) {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Bind context to auth %v", context)
	authenticate := new(adapter.AuthenticateDTO)
	if err := context.Bind(authenticate); err != nil {
		logger.Warnf("Could not bind auth, %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}
	logger.Debugf("Auth bind %v", authenticate)
	if err := context.Validate(authenticate); err != nil {
		logger.Warnf("Could not validate auth, %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		logger.Warnf("Error while binding userId: %v", err)
		return uuid.Nil, nil, echo.ErrBadRequest
	}

	return userId, authenticate, nil
}

func checkUserId(context echo.Context, userId uuid.UUID) error {
	logger := context.Get(loggerKey).(*log.Entry)
	claims, ok := context.Get(adapter.ClaimName).(adapter.Claims)

	if !ok {
		logger.Warnf("Could not map Claims")
		return echo.ErrUnauthorized
	}

	if userId != claims.UserId {
		logger.Warnf("User %v is not allowed to get token for user %v", userId, claims.UserId)
		return echo.ErrUnauthorized
	}

	return nil
}

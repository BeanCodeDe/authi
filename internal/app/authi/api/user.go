package api

import (
	"errors"
	"net/http"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/BeanCodeDe/authi/pkg/authmiddleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const (
	UserRootPath = "/user"
	userIdParam  = "userId"
)

type (
	authenticate struct {
		Password string `json:"password" validate:"required"`
	}

	tokenResponseDTO struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
	}
)

func InitUserInterface(group *echo.Group) {
	core.Init()
	group.POST("", createUserId)
	group.POST("/:"+userIdParam+"/login", login)
	group.PUT("/:"+userIdParam, create)
	group.PATCH("/:"+userIdParam+"/refresh", refreshToken, authmiddleware.CheckToken)
}

func createUserId(context echo.Context) error {
	log.Debug("Create User")
	return context.String(http.StatusOK, uuid.NewString())
}

func create(context echo.Context) error {
	log.Debugf("Create user")
	userCore, err := bindUser(context)
	if err != nil {
		return err
	}

	if err := userCore.Create(); err != nil {
		if errors.Is(err, errormessages.UserAlreadyExists) {
			log.Warnf("User with id %s already exists", userCore.GetId())
			return context.NoContent(http.StatusConflict)
		}
		return err
	}

	log.Debugf("Created user with id %s", userCore.GetId())
	return context.NoContent(http.StatusCreated)
}

func refreshToken(context echo.Context) error {
	log.Debugf("Refresh token")

	refreshToken := context.Request().Header.Get(authadapter.RefreshTokenHeaderName)
	claims, ok := context.Get(authadapter.ClaimName).(authadapter.Claims)

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	if userId != claims.UserId {
		log.Warnf("User %v is not allowed to get token for user %v", userId, claims.UserId)
		return echo.ErrUnauthorized
	}

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
	userCore, err := bindUser(context)
	if err != nil {
		return err
	}

	tokenCore, err := userCore.Login()
	if err != nil {
		log.Warnf("Error while logging in user %v: %v", userCore, err)
		return echo.ErrUnauthorized
	}

	log.Debugf("Logged in user %s", userCore.GetId())
	return context.JSON(http.StatusOK, mapToTokenResponseDTO(tokenCore))
}

func bindUser(context echo.Context) (*core.UserCore, error) {
	log.Debugf("Bind context to user %v", context)
	authenticate := new(authenticate)
	if err := context.Bind(authenticate); err != nil {
		log.Warnf("Could not bind user, %v", err)
		return nil, echo.ErrBadRequest
	}
	log.Debugf("User bind %v", authenticate)
	if err := context.Validate(authenticate); err != nil {
		log.Warnf("Could not validate user, %v", err)
		return nil, echo.ErrBadRequest
	}

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return nil, echo.ErrBadRequest
	}
	userCore := &core.UserCore{ID: userId, Password: authenticate.Password}
	return userCore, nil
}

func mapToTokenResponseDTO(tokenCore *core.TokenCore) *tokenResponseDTO {
	return &tokenResponseDTO{AccessToken: tokenCore.AccessToken, ExpiresIn: tokenCore.ExpiresIn, RefreshToken: tokenCore.RefreshToken, RefreshExpiresIn: tokenCore.RefreshExpiresIn}
}

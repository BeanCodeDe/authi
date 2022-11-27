package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/BeanCodeDe/authi/pkg/authmiddleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type (
	UserApi struct {
		facade core.Facade
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewUserApi(auth authadapter.Auth) (*UserApi, error) {
	userFacade, err := core.NewUserFacade()
	if err != nil {
		return nil, fmt.Errorf("error while initializing user facade: %v", err)
	}

	authMiddleware := authmiddleware.NewAuthmiddleware(auth)
	api := &UserApi{userFacade}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{validator: validator.New()}

	userGroup := e.Group(userRootPath)
	userGroup.POST("", api.CreateUserId)
	userGroup.POST("/:"+userIdParam+userLoginPath, api.LoginUser)
	userGroup.PUT("/:"+userIdParam, api.CreateUser)
	userGroup.PATCH("/:"+userIdParam+userRefreshPath, api.RefreshToken, authMiddleware.CheckToken)

	e.Logger.Fatal(e.Start(":1203"))

	return api, nil
}

func (userApi *UserApi) CreateUserId(context echo.Context) error {
	log.Debug("Create User")
	return context.String(http.StatusCreated, uuid.NewString())
}

func (userApi *UserApi) CreateUser(context echo.Context) error {
	log.Debugf("Create user")
	userId, authenticate, err := userApi.bindAuthenticate(context)
	if err != nil {
		return err
	}

	if err := userApi.facade.CreateUser(userId, authenticate); err != nil {
		if errors.Is(err, errormessages.ErrUserAlreadyExists) {
			log.Warnf("User with id %s already exists", userId)
			return echo.NewHTTPError(http.StatusConflict)
		}
		log.Warnf("Error while creating user: %v", err)
		return echo.ErrInternalServerError
	}

	log.Debugf("Created user with id %s", userId)
	return context.NoContent(http.StatusCreated)
}

func (userApi *UserApi) LoginUser(context echo.Context) error {
	log.Debugf("Login some user")
	userId, authenticate, err := userApi.bindAuthenticate(context)
	if err != nil {
		return err
	}

	token, err := userApi.facade.LoginUser(userId, authenticate)
	if err != nil {
		log.Warnf("Error while logging in user %v: %v", userId, err)
		return echo.ErrUnauthorized
	}

	log.Debugf("Logged in user %s", userId)
	return context.JSON(http.StatusOK, token)
}

func (userApi *UserApi) RefreshToken(context echo.Context) error {
	log.Debugf("Refresh token")

	refreshToken := context.Request().Header.Get(authadapter.RefreshTokenHeaderName)
	claims, ok := context.Get(authadapter.ClaimName).(authadapter.Claims)

	if !ok {
		log.Errorf("Got data of wrong type: %v", context.Get(authadapter.ClaimName))
		return echo.ErrUnauthorized
	}

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	if userId != claims.UserId {
		log.Warnf("User %v is not allowed to get token for user %v", userId, claims.UserId)
		return echo.ErrUnauthorized
	}

	token, err := userApi.facade.RefreshToken(claims.UserId, refreshToken)
	if err != nil {
		log.Errorf("Something went wrong while creating Token: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("Refresh token for user %s updated", claims.UserId)
	return context.JSON(http.StatusOK, token)
}

func (userApi *UserApi) bindAuthenticate(context echo.Context) (uuid.UUID, *core.AuthenticateDTO, error) {
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

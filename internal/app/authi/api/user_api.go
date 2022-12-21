package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	echoMiddleware "github.com/BeanCodeDe/authi/pkg/middleware"
	"github.com/BeanCodeDe/authi/pkg/parser"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var (
	userIdParam = "userId"
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

func MiddlewareExample() {
	//Initialize parser to validate Tokens
	tokenParser, err := parser.NewJWTParser()

	//Checking if an error occurred while loading jwt parser
	if err != nil {
		panic(err)
	}

	//Initialize middleware
	echoMiddleware := echoMiddleware.NewEchoMiddleware(tokenParser)

	//Initialize echo
	e := echo.New()

	//Secure endpoint with method `echoMiddleware.CheckToken`
	e.GET(
		"/someEndpoint",
		func(c echo.Context) error { return c.NoContent(201) },
		echoMiddleware.CheckToken,
	)
}

func NewUserApi(parser parser.Parser) (*UserApi, error) {
	userFacade, err := core.NewUserFacade()
	if err != nil {
		return nil, fmt.Errorf("error while initializing user facade: %v", err)
	}

	echoMiddleware := echoMiddleware.NewEchoMiddleware(parser)
	api := &UserApi{userFacade}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{validator: validator.New()}

	userGroup := e.Group(adapter.AuthiRootPath)
	userGroup.POST("", api.CreateUserId)
	userGroup.POST("/:"+userIdParam+adapter.AuthiLoginPath, api.LoginUser)
	userGroup.PUT("/:"+userIdParam, api.CreateUser)
	userGroup.PATCH("/:"+userIdParam+adapter.AuthiRefreshPath, api.RefreshToken, echoMiddleware.CheckToken)
	userGroup.PATCH("/:"+userIdParam, api.UpdatePassword, echoMiddleware.CheckToken)
	userGroup.DELETE("/:"+userIdParam, api.DeleteUser, echoMiddleware.CheckToken)

	address := util.GetEnvWithFallback("ADDRESS", "0.0.0.0")
	port, err := util.GetEnvIntWithFallback("PORT", 1203)
	if err != nil {
		return nil, fmt.Errorf("error while loading port from environment variable: %w", err)
	}
	url := fmt.Sprintf("%s:%d", address, port)
	e.Logger.Fatal(e.Start(url))

	return api, nil
}

func (userApi *UserApi) CreateUserId(context echo.Context) error {
	log.Debug("Create User Id")
	return context.String(http.StatusCreated, uuid.NewString())
}

func (userApi *UserApi) CreateUser(context echo.Context) error {
	log.Debugf("Create User")
	userId, authenticate, err := bindAuthenticate(context)
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
	userId, authenticate, err := bindAuthenticate(context)
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

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	refreshToken := context.Request().Header.Get(adapter.RefreshTokenHeaderName)

	token, err := userApi.facade.RefreshToken(userId, refreshToken)
	if err != nil {
		log.Errorf("Something went wrong while creating Token: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("Refresh token for user %s updated", userId)
	return context.JSON(http.StatusOK, token)
}

func (userApi *UserApi) UpdatePassword(context echo.Context) error {
	log.Debugf("Update password")

	userId, authenticate, err := bindAuthenticate(context)
	if err != nil {
		return err
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	err = userApi.facade.UpdatePassword(userId, authenticate)
	if err != nil {
		log.Errorf("Something went wrong while updating password: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("Password for user %s updated", userId)
	return context.NoContent(http.StatusNoContent)
}

func (userApi *UserApi) DeleteUser(context echo.Context) error {
	log.Debugf("Delete password")

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	err = userApi.facade.DeleteUser(userId)
	if err != nil {
		log.Errorf("Something went wrong while deleting user: %v", err)
		return echo.ErrUnauthorized
	}
	log.Debugf("User %s deleted", userId)
	return context.NoContent(http.StatusNoContent)
}

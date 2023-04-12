package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	echoMiddleware "github.com/BeanCodeDe/authi/pkg/middleware"
	"github.com/BeanCodeDe/authi/pkg/parser"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/go-playground/validator.v9"
)

var (
	userIdParam = "userId"
)

const (
	loggerKey = "logger"
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

func NewUserApi(parser parser.Parser) (*UserApi, error) {
	userFacade, err := core.NewUserFacade()
	if err != nil {
		return nil, fmt.Errorf("error while initializing user facade: %v", err)
	}

	echoMiddleware := echoMiddleware.NewEchoMiddleware(parser)
	api := &UserApi{userFacade}

	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(middleware.CORS(), setLoggerMiddleware, middleware.Recover())

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			userId := ctx.Param(userIdParam)
			if userId != "" {
				return userId, nil
			}
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

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
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debug("Create User Id")
	return context.String(http.StatusCreated, uuid.NewString())
}

func (userApi *UserApi) CreateUser(context echo.Context) error {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Create User")
	userId, authenticate, err := bindAuthenticate(context)
	if err != nil {
		return err
	}

	if err := userApi.facade.CreateUser(userId, authenticate.Password, false); err != nil {
		logger.Warnf("Error while creating user: %v", err)
		return echo.ErrUnauthorized
	}

	logger.Debugf("Created user with id %s", userId)
	return context.NoContent(http.StatusCreated)
}

func (userApi *UserApi) LoginUser(context echo.Context) error {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Login some user")
	userId, authenticate, err := bindAuthenticate(context)
	if err != nil {
		return err
	}

	token, err := userApi.facade.LoginUser(userId, authenticate.Password)
	if err != nil {
		logger.Warnf("Error while logging in user %v: %v", userId, err)
		return echo.ErrUnauthorized
	}

	logger.Debugf("Logged in user %s", userId)
	return context.JSON(http.StatusOK, token)
}

func (userApi *UserApi) RefreshToken(context echo.Context) error {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Refresh token")

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		logger.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	refreshToken := context.Request().Header.Get(adapter.RefreshTokenHeaderName)

	token, err := userApi.facade.RefreshToken(userId, refreshToken)
	if err != nil {
		logger.Errorf("Something went wrong while creating Token: %v", err)
		return echo.ErrUnauthorized
	}
	logger.Debugf("Refresh token for user %s updated", userId)
	return context.JSON(http.StatusOK, token)
}

func (userApi *UserApi) UpdatePassword(context echo.Context) error {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Update password")

	userId, authenticate, err := bindAuthenticate(context)
	if err != nil {
		return err
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	err = userApi.facade.UpdatePassword(userId, authenticate.Password)
	if err != nil {
		logger.Errorf("Something went wrong while updating password: %v", err)
		return echo.ErrUnauthorized
	}
	logger.Debugf("Password for user %s updated", userId)
	return context.NoContent(http.StatusNoContent)
}

func (userApi *UserApi) DeleteUser(context echo.Context) error {
	logger := context.Get(loggerKey).(*log.Entry)
	logger.Debugf("Delete password")

	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		logger.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}

	err = checkUserId(context, userId)
	if err != nil {
		return err
	}

	err = userApi.facade.DeleteUser(userId)
	if err != nil {
		logger.Errorf("Something went wrong while deleting user: %v", err)
		return echo.ErrUnauthorized
	}
	logger.Debugf("User %s deleted", userId)
	return context.NoContent(http.StatusNoContent)
}

func setLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		correlationId := c.Request().Header.Get(CorrelationIdHeader)
		_, err := uuid.Parse(correlationId)
		if err != nil {
			log.Warn("Correlation id is not from format uuid. Set default correlation id.")
			correlationId = "WRONG FORMAT"
		}
		logger := log.WithFields(log.Fields{
			CorrelationIdHeader: correlationId,
		})

		c.Set(loggerKey, logger)
		return next(c)
	}
}

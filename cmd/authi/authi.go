package main

import (
	"github.com/BeanCodeDe/authi/internal/app/authi/api"
	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/BeanCodeDe/authi/pkg/authmiddleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

const (
	userRootPath = "/user"
	userIdParam  = "userId"
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	setLogLevel(config.LogLevel)
	log.Info("Start Server")

	authAdapter, err := authadapter.NewAuthAdapter()
	if err != nil {
		log.Fatalf("Error while initializing auth adapter: %v", err)
	}

	userApi, err := api.NewUserApi(authAdapter)
	if err != nil {
		log.Fatalf("Error while initializing user api: %v", err)
	}

	authMiddleware := authmiddleware.NewAuthmiddleware(authAdapter)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{validator: validator.New()}

	userGroup := e.Group(userRootPath)
	userGroup.POST("", userApi.CreateUserId)
	userGroup.POST("/:"+userIdParam+"/login", userApi.LoginUser)
	userGroup.PUT("/:"+userIdParam, userApi.CreateUser)
	userGroup.PATCH("/:"+userIdParam+"/refresh", userApi.RefreshToken, authMiddleware.CheckToken)

	e.Logger.Fatal(e.Start(":1203"))
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	}
}

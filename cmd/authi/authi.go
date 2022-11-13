package main

import (
	"github.com/BeanCodeDe/authi/internal/app/authi/api"
	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	defer handleExit()
	setLogLevel(config.LogLevel)
	log.Info("Start Server")
	db.Init()
	e := echo.New()
	e.HTTPErrorHandler = api.CustomHTTPErrorHandler
	e.Validator = &CustomValidator{validator: validator.New()}
	userGroup := e.Group(api.UserRootPath)
	api.InitUserInterface(userGroup)
	authGroup := e.Group(api.AuthRootPath)
	api.InitAuthInterface(authGroup)
	e.Logger.Fatal(e.Start(":1323"))
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

func handleExit() {
	db.Close()
}

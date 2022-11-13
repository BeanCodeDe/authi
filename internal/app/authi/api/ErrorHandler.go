package api

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	log.Error("An unhandlet error accurd: %v", err)
	c.NoContent(echo.ErrInternalServerError.Code)
}

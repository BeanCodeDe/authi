package api

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	log.Warnf("An Error accurd: %v", err)

	echoError, ok := err.(*echo.HTTPError)
	if ok {
		c.String(echoError.Code, "")
		return
	}

	c.String(echo.ErrUnauthorized.Code, "")
}

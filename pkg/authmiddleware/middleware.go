package authmiddleware

import (
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type (
	Authmiddleware struct {
		auth authadapter.Auth
	}

	Middleware interface {
		CheckToken(next echo.HandlerFunc) echo.HandlerFunc
	}
)

func NewAuthmiddleware(auth authadapter.Auth) *Authmiddleware {
	return &Authmiddleware{auth}
}

func (authmiddleware *Authmiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authadapter.AuthorizationHeaderName)

		claims, err := authmiddleware.auth.ParseToken(authHeader)
		if err != nil {
			log.Warnf("error while parsing token %v", err)
			return echo.ErrUnauthorized
		}

		c.Set(authadapter.ClaimName, *claims)
		return next(c)
	}
}

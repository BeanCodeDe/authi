package authmiddleware

import (
	"fmt"
	"time"

	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/labstack/echo/v4"
)

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authadapter.AuthorizationHeaderName)
		if authHeader == "" {
			return fmt.Errorf("no auth Header found")
		}

		claims, err := authadapter.ParseToken(authHeader)
		if err != nil {
			return fmt.Errorf("error while parsing token: %v", err)
		}

		var token string

		if time.Now().Add(1 * time.Minute).After(time.Unix(claims.ExpiresAt, 0)) {
			token, err = authadapter.CreateJWTToken(authHeader)
			if err != nil {
				return fmt.Errorf("error while creating token: %v", err)
			}
			c.Response().Header().Set(authadapter.AuthorizationHeaderName, token)
		} else {
			token = authHeader
		}
		c.Set(authadapter.ClaimName, *claims)
		c.Response().Header().Set(authadapter.AuthorizationHeaderName, token)
		return next(c)
	}
}

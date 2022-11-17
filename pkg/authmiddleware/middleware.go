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

		claims, err := authadapter.ParseToken(authHeader)
		if err != nil {
			fmt.Printf("error while parsing token %v", err)
			//error while parsing token
			return echo.ErrBadRequest
		}

		if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
			fmt.Printf("token expired")
			//token expired
			return echo.ErrUnauthorized
		}

		c.Set(authadapter.ClaimName, *claims)
		return next(c)
	}
}

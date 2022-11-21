package api

import "github.com/labstack/echo/v4"

type (
	Api interface {
		CreateUserId(context echo.Context) error
		CreateUser(context echo.Context) error
		RefreshToken(context echo.Context) error
		LoginUser(context echo.Context) error
	}
)

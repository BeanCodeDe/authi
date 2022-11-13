package authErr

import "github.com/labstack/echo/v4"

type CustomError struct {
	Msg      string
	HttpCode int
}

func (m *CustomError) Error() string {
	return m.Msg
}

var (
	UserNotFoundError  = &CustomError{"user not found", echo.ErrNotFound.Code}
	WrongAuthDataError = &CustomError{"wrong auth data", echo.ErrUnauthorized.Code}
	UnknownError       = &CustomError{"unknown error", echo.ErrInternalServerError.Code}
)

package examples

import (
	"github.com/BeanCodeDe/authi/pkg/middleware"
	"github.com/BeanCodeDe/authi/pkg/parser"
	"github.com/labstack/echo/v4"
)

func MiddlewareExample() {
	//Initialize parser to validate Tokens
	tokenParser, err := parser.NewJWTParser()

	//Checking if an error occurred while loading jwt parser
	if err != nil {
		panic(err)
	}

	//Initialize middleware
	echoMiddleware := middleware.NewEchoMiddleware(tokenParser)

	//Initialize echo
	e := echo.New()

	//Secure endpoint with method `echoMiddleware.CheckToken`
	e.GET(
		"/someEndpoint",
		func(c echo.Context) error { return c.NoContent(201) },
		echoMiddleware.CheckToken,
	)
}

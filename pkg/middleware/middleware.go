package middleware

import (
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/BeanCodeDe/authi/pkg/parser"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type (
	//Implementation of middleware interface
	EchoMiddleware struct {
		auth        adapter.AuthAdapter
		tokenParser parser.Parser
	}

	//Interface to check token
	Middleware interface {
		CheckToken(next echo.HandlerFunc) echo.HandlerFunc
	}
)

// Constructor to creat new EchoMiddleware
func NewEchoMiddleware(auth adapter.AuthAdapter, tokenParser parser.Parser) Middleware {
	return &EchoMiddleware{auth, tokenParser}
}

// Check if incoming token is valid and set claim in context with key claim
func (middleware *EchoMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(adapter.AuthorizationHeaderName)

		claims, err := middleware.tokenParser.ParseToken(authHeader)
		if err != nil {
			log.Warnf("error while parsing token %v", err)
			return echo.ErrUnauthorized
		}

		c.Set(adapter.ClaimName, *claims)
		return next(c)
	}
}

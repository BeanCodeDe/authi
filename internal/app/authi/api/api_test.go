package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

var (
	errSome                       = errors.New("some error from facade")
	successfullyTokenResponse     = []*core.AuthenticateResponse{{TokenResponse: &adapter.TokenResponseDTO{AccessToken: "some_access_token", ExpiresIn: 1, RefreshToken: "some_refresh_token", RefreshExpiresIn: 2}, Err: nil}}
	errorTokenResponse            = []*core.AuthenticateResponse{{TokenResponse: nil, Err: errSome}}
	userId                        = uuid.New()
	wrongUUID                     = "xyz"
	password                      = "some_password"
	refreshToken                  = "some_refresh_token"
	authenticationUserJson        = fmt.Sprintf(`{"password":"%s"}`, password)
	authenticateObject            = &adapter.AuthenticateDTO{Password: password}
	authenticationUserInvalidJson = `{"password":""}`
	claimUser                     = adapter.Claims{UserId: userId}
	wrongClaimFormat              = &UserApi{}
)

// bindAuthenticate tests

func TestBindAuthenticate_Successfully(t *testing.T) {
	// Prep
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())

	// Exec
	returnedUserId, returnedAuthenticate, returnedErr := bindAuthenticate(c)

	// Assertions
	assert.Equal(t, userId, returnedUserId)
	assert.Equal(t, authenticateObject, returnedAuthenticate)
	assert.Nil(t, returnedErr)
}

func TestBindAuthenticate_CouldNotBind(t *testing.T) {
	// Prep
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	c := e.NewContext(req, nil)
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())

	// Exec
	returnedUserId, returnedAuthenticate, returnedErr := bindAuthenticate(c)

	// Assertions
	assert.Equal(t, uuid.Nil, returnedUserId)
	assert.Nil(t, returnedAuthenticate)
	assert.Equal(t, echo.ErrBadRequest, returnedErr)
}

func TestBindAuthenticate_ValidateError(t *testing.T) {
	// Prep
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserInvalidJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())

	// Exec
	returnedUserId, returnedAuthenticate, returnedErr := bindAuthenticate(c)

	// Assertions
	assert.Equal(t, uuid.Nil, returnedUserId)
	assert.Nil(t, returnedAuthenticate)
	assert.Equal(t, echo.ErrBadRequest, returnedErr)
}

func TestBindAuthenticate_ParseError(t *testing.T) {
	// Prep
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(wrongUUID)

	// Exec
	returnedUserId, returnedAuthenticate, returnedErr := bindAuthenticate(c)

	// Assertions
	assert.Equal(t, uuid.Nil, returnedUserId)
	assert.Nil(t, returnedAuthenticate)
	assert.Equal(t, echo.ErrBadRequest, returnedErr)
}

// checkUserId tests

func TestCheckUserId_Successfully(t *testing.T) {
	// Prep
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set(adapter.ClaimName, claimUser)

	// Exec
	returnedErr := checkUserId(c, userId)

	// Assertions
	assert.Nil(t, returnedErr)
}

func TestCheckUserId_CouldNotMapClaim(t *testing.T) {
	// Prep
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set(adapter.ClaimName, wrongClaimFormat)

	// Exec
	returnedErr := checkUserId(c, userId)

	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, returnedErr)
}

func TestCheckUserId_ClaimDoesNotMatchUserId(t *testing.T) {
	// Prep
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set(adapter.ClaimName, claimUser)

	// Exec
	returnedErr := checkUserId(c, uuid.New())

	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, returnedErr)
}

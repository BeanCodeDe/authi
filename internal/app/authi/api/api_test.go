package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

type (
	authenticateRecord struct {
		userId       uuid.UUID
		authenticate *core.AuthenticateDTO
	}

	refreshTokenRecord struct {
		userId uuid.UUID
		token  string
	}

	deleteUserRecord struct {
		userId uuid.UUID
	}

	authenticateReturn struct {
		tokenResponse *core.TokenResponseDTO
		err           error
	}

	facadeMock struct {
		createUserRecordArray     []*authenticateRecord
		loginUserRecordArray      []*authenticateRecord
		refreshTokenRecordArray   []*refreshTokenRecord
		updatePasswordRecordArray []*authenticateRecord
		deleteUserRecordArray     []*deleteUserRecord
		createUserReturn          []error
		loginUserReturn           []*authenticateReturn
		refreshTokenReturn        []*authenticateReturn
		updatePasswordReturn      []error
		deleteUserReturn          []error
	}
)

var (
	errSome                       = errors.New("some error from facade")
	successfullyTokenResponse     = []*authenticateReturn{{tokenResponse: &core.TokenResponseDTO{AccessToken: "some_access_token", ExpiresIn: 1, RefreshToken: "some_refresh_token", RefreshExpiresIn: 2}, err: nil}}
	errorTokenResponse            = []*authenticateReturn{{tokenResponse: nil, err: errSome}}
	userId                        = uuid.New()
	wrongUUID                     = "xyz"
	password                      = "some_password"
	refreshToken                  = "some_refresh_token"
	authenticationUserJson        = fmt.Sprintf(`{"password":"%s"}`, password)
	authenticateObject            = &core.AuthenticateDTO{Password: password}
	authenticationUserWrongJson   = "xyz"
	authenticationUserInvalidJson = `{"password":""}`
	claimUser                     = authadapter.Claims{UserId: userId}
	wrongClaimFormat              = &UserApi{}
)

// bindAuthenticate tests

func TestBindAuthenticate_Successfully(t *testing.T) {
	// Prep
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(userRootPath + "/:" + userIdParam)
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
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	c := e.NewContext(req, nil)
	c.SetPath(userRootPath + "/:" + userIdParam)
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
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserInvalidJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(userRootPath + "/:" + userIdParam)
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
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, nil)
	c.SetPath(userRootPath + "/:" + userIdParam)
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
	c.Set(authadapter.ClaimName, claimUser)

	// Exec
	returnedErr := checkUserId(c, userId)

	// Assertions
	assert.Nil(t, returnedErr)
}

func TestCheckUserId_CouldNotMapClaim(t *testing.T) {
	// Prep
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set(authadapter.ClaimName, wrongClaimFormat)

	// Exec
	returnedErr := checkUserId(c, userId)

	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, returnedErr)
}

func TestCheckUserId_ClaimDoesNotMatchUserId(t *testing.T) {
	// Prep
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set(authadapter.ClaimName, claimUser)

	// Exec
	returnedErr := checkUserId(c, uuid.New())

	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, returnedErr)
}

//Mock methods

func (facadeMock *facadeMock) CreateUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	createUserRecord := &authenticateRecord{userId, authenticate}
	facadeMock.createUserRecordArray = append(facadeMock.createUserRecordArray, createUserRecord)
	return facadeMock.createUserReturn[len(facadeMock.createUserRecordArray)-1]
}

func (facadeMock *facadeMock) LoginUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) (*core.TokenResponseDTO, error) {
	loginUserRecord := &authenticateRecord{userId, authenticate}
	facadeMock.loginUserRecordArray = append(facadeMock.loginUserRecordArray, loginUserRecord)
	loginReturn := facadeMock.loginUserReturn[len(facadeMock.loginUserRecordArray)-1]
	return loginReturn.tokenResponse, loginReturn.err
}

func (facadeMock *facadeMock) RefreshToken(userId uuid.UUID, refreshToken string) (*core.TokenResponseDTO, error) {
	refreshTokenRecord := &refreshTokenRecord{userId, refreshToken}
	facadeMock.refreshTokenRecordArray = append(facadeMock.refreshTokenRecordArray, refreshTokenRecord)
	loginReturn := facadeMock.refreshTokenReturn[len(facadeMock.refreshTokenRecordArray)-1]
	return loginReturn.tokenResponse, loginReturn.err
}

func (facadeMock *facadeMock) UpdatePassword(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	updatePasswordRecord := &authenticateRecord{userId, authenticate}
	facadeMock.updatePasswordRecordArray = append(facadeMock.updatePasswordRecordArray, updatePasswordRecord)
	return facadeMock.updatePasswordReturn[len(facadeMock.updatePasswordRecordArray)-1]
}

func (facadeMock *facadeMock) DeleteUser(userId uuid.UUID) error {
	deleteUserRecord := &deleteUserRecord{userId}
	facadeMock.deleteUserRecordArray = append(facadeMock.deleteUserRecordArray, deleteUserRecord)
	return facadeMock.deleteUserReturn[len(facadeMock.deleteUserRecordArray)-1]
}

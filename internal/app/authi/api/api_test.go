package api

import (
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
	createUserRecord struct {
		userId       uuid.UUID
		authenticate *core.AuthenticateDTO
	}

	loginUserRecord struct {
		userId       uuid.UUID
		authenticate *core.AuthenticateDTO
	}

	refreshUserRecord struct {
		userId uuid.UUID
		token  string
	}

	facadeMock struct {
		createUserRecordArray  []*createUserRecord
		loginUserRecordArray   []*loginUserRecord
		refreshUserRecordArray []*refreshUserRecord
	}
)

var (
	successfullyTokenResponse = &core.TokenResponseDTO{AccessToken: "some_access_token", ExpiresIn: 1, RefreshToken: "some_refresh_token", RefreshExpiresIn: 2}
	userId                    = uuid.New()
	password                  = "some_password"
	refreshToken              = "some_refresh_token"
	authenticationUserJson    = fmt.Sprintf(`{"password":"%s"}`, password)
	claimUser                 = authadapter.Claims{UserId: userId}
)

func TestCreateUserId_Successfully(t *testing.T) {
	userApi := &UserApi{&facadeMock{}}
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, userRootPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, userApi.CreateUserId(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateUser_Successfully(t *testing.T) {
	facade := &facadeMock{}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.NoError(t, userApi.CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 1, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
		assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)
	}
}

func (facadeMock *facadeMock) CreateUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	createUserRecord := &createUserRecord{userId, authenticate}
	facadeMock.createUserRecordArray = append(facadeMock.createUserRecordArray, createUserRecord)
	return nil
}

func TestLoginUser_Successfully(t *testing.T) {
	facade := &facadeMock{}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userLoginPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.NoError(t, userApi.LoginUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 1, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.loginUserRecordArray[0].userId)
		assert.Equal(t, password, facade.loginUserRecordArray[0].authenticate.Password)
		assert.Equal(t, "{\"access_token\":\"some_access_token\",\"expires_in\":1,\"refresh_token\":\"some_refresh_token\",\"refresh_expires_in\":2}\n", rec.Body.String())
	}
}

func (facadeMock *facadeMock) LoginUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) (*core.TokenResponseDTO, error) {
	loginUserRecord := &loginUserRecord{userId, authenticate}
	facadeMock.loginUserRecordArray = append(facadeMock.loginUserRecordArray, loginUserRecord)
	return successfullyTokenResponse, nil
}

func TestRefreshToken_Successfully(t *testing.T) {
	facade := &facadeMock{}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(authadapter.RefreshTokenHeaderName, refreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userRefreshPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Assertions
	if assert.NoError(t, userApi.RefreshToken(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 1, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.refreshUserRecordArray[0].userId)
		assert.Equal(t, refreshToken, facade.refreshUserRecordArray[0].token)
	}
}

func (facadeMock *facadeMock) RefreshToken(userId uuid.UUID, refreshToken string) (*core.TokenResponseDTO, error) {
	refreshUserRecord := &refreshUserRecord{userId, refreshToken}
	facadeMock.refreshUserRecordArray = append(facadeMock.refreshUserRecordArray, refreshUserRecord)
	return successfullyTokenResponse, nil
}

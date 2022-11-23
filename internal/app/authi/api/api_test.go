package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
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
		t                      *testing.T
		createUserRecordArray  []*createUserRecord
		loginUserRecordArray   []*loginUserRecord
		refreshUserRecordArray []*refreshUserRecord
	}

	facadeDuplicateMock struct {
		t                      *testing.T
		createUserRecordArray  []*createUserRecord
		loginUserRecordArray   []*loginUserRecord
		refreshUserRecordArray []*refreshUserRecord
	}

	facadeErrorMock struct {
		t                      *testing.T
		createUserRecordArray  []*createUserRecord
		loginUserRecordArray   []*loginUserRecord
		refreshUserRecordArray []*refreshUserRecord
	}
)

var (
	successfullyTokenResponse     = &core.TokenResponseDTO{AccessToken: "some_access_token", ExpiresIn: 1, RefreshToken: "some_refresh_token", RefreshExpiresIn: 2}
	userId                        = uuid.New()
	wrongUUID                     = "xyz"
	password                      = "some_password"
	refreshToken                  = "some_refresh_token"
	authenticationUserJson        = fmt.Sprintf(`{"password":"%s"}`, password)
	authenticationUserWrongJson   = "xyz"
	authenticationUserInvalidJson = `{"password":""}`
	claimUser                     = authadapter.Claims{UserId: userId}
	wrongClaimFormat              = &UserApi{}
)

// CreateUserId Tests

func TestCreateUserId_Successfully(t *testing.T) {
	userApi := &UserApi{&facadeMock{t: t}}
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

// CreateUser Tests

func TestCreateUser_Successfully(t *testing.T) {
	facade := &facadeMock{t: t}
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

func TestCreateUser_BindAuth_Error(t *testing.T) {
	facade := &facadeMock{t: t}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserWrongJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestCreateUser_ValidateAuth_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserInvalidJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestCreateUser_ParseUserIdParam_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
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
	c.SetParamValues(wrongUUID)
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func (facadeMock *facadeMock) CreateUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	createUserRecord := &createUserRecord{userId, authenticate}
	facadeMock.createUserRecordArray = append(facadeMock.createUserRecordArray, createUserRecord)
	return nil
}

func TestCreateUser_CreateUser_ErrUserAlreadyExists(t *testing.T) {
	facade := &facadeDuplicateMock{t: t}
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
	if assert.Equal(t, echo.NewHTTPError(http.StatusConflict), userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 1, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
		assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)
	}
}

func (facadeMock *facadeDuplicateMock) CreateUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	createUserRecord := &createUserRecord{userId, authenticate}
	facadeMock.createUserRecordArray = append(facadeMock.createUserRecordArray, createUserRecord)
	return errormessages.ErrUserAlreadyExists
}

func TestCreateUser_CreateUser_InternalServerError(t *testing.T) {
	facade := &facadeErrorMock{t: t}
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
	if assert.Equal(t, echo.ErrInternalServerError, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 1, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
		assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)
	}
}

func (facadeMock *facadeErrorMock) CreateUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) error {
	createUserRecord := &createUserRecord{userId, authenticate}
	facadeMock.createUserRecordArray = append(facadeMock.createUserRecordArray, createUserRecord)
	return echo.ErrInternalServerError
}

// LoginUser Tests

func TestLoginUser_Successfully(t *testing.T) {
	facade := &facadeMock{t: t}
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

func TestLoginUser_BindAuth_Error(t *testing.T) {
	facade := &facadeMock{t: t}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserWrongJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userLoginPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestLoginUser_ValidateAuth_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, strings.NewReader(authenticationUserInvalidJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userLoginPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestLoginUser_ParseUserIdParam_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
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
	c.SetParamValues(wrongUUID)
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.CreateUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func (facadeMock *facadeMock) LoginUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) (*core.TokenResponseDTO, error) {
	loginUserRecord := &loginUserRecord{userId, authenticate}
	facadeMock.loginUserRecordArray = append(facadeMock.loginUserRecordArray, loginUserRecord)
	return successfullyTokenResponse, nil
}

func (facadeMock *facadeDuplicateMock) LoginUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) (*core.TokenResponseDTO, error) {
	assert.Fail(facadeMock.t, "Login user of facade duplicate mock shouldn't be access")
	return nil, nil
}

func TestLoginUser__LoginUser_ErrUnauthorized(t *testing.T) {
	facade := &facadeErrorMock{t: t}
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
	if assert.Equal(t, echo.ErrUnauthorized, userApi.LoginUser(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 1, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.loginUserRecordArray[0].userId)
		assert.Equal(t, password, facade.loginUserRecordArray[0].authenticate.Password)
	}
}

func (facadeMock *facadeErrorMock) LoginUser(userId uuid.UUID, authenticate *core.AuthenticateDTO) (*core.TokenResponseDTO, error) {
	loginUserRecord := &loginUserRecord{userId, authenticate}
	facadeMock.loginUserRecordArray = append(facadeMock.loginUserRecordArray, loginUserRecord)
	return nil, echo.ErrInternalServerError
}

// RefreshToken Tests

func TestRefreshToken_Successfully(t *testing.T) {
	facade := &facadeMock{t: t}
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

func TestRefreshToken_ParseClaim_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
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
	c.Set(authadapter.ClaimName, wrongClaimFormat)
	// Assertions
	if assert.Equal(t, echo.ErrUnauthorized, userApi.RefreshToken(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestRefreshToken_ParseUserIdParam_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
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
	c.SetParamValues(wrongUUID)
	c.Set(authadapter.ClaimName, claimUser)
	// Assertions
	if assert.Equal(t, echo.ErrBadRequest, userApi.RefreshToken(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func TestRefreshToken_CheckUserId_Invalid(t *testing.T) {
	facade := &facadeMock{t: t}
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
	c.SetParamValues(uuid.NewString())
	c.Set(authadapter.ClaimName, claimUser)
	// Assertions
	if assert.Equal(t, echo.ErrUnauthorized, userApi.RefreshToken(c)) {
		assert.Equal(t, 0, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
	}
}

func (facadeMock *facadeMock) RefreshToken(userId uuid.UUID, refreshToken string) (*core.TokenResponseDTO, error) {
	refreshUserRecord := &refreshUserRecord{userId, refreshToken}
	facadeMock.refreshUserRecordArray = append(facadeMock.refreshUserRecordArray, refreshUserRecord)
	return successfullyTokenResponse, nil
}

func (facadeMock *facadeDuplicateMock) RefreshToken(userId uuid.UUID, refreshToken string) (*core.TokenResponseDTO, error) {
	assert.Fail(facadeMock.t, "Login user of facade duplicate mock shouldn't be access")
	return nil, nil
}

func TestRefreshToken_RefreshToken_ErrUnauthorized(t *testing.T) {
	facade := &facadeErrorMock{t: t}
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
	if assert.Equal(t, echo.ErrUnauthorized, userApi.RefreshToken(c)) {
		assert.Equal(t, 1, len(facade.refreshUserRecordArray))
		assert.Equal(t, 0, len(facade.loginUserRecordArray))
		assert.Equal(t, 0, len(facade.createUserRecordArray))
		assert.Equal(t, userId, facade.refreshUserRecordArray[0].userId)
		assert.Equal(t, refreshToken, facade.refreshUserRecordArray[0].token)
	}
}

func (facadeMock *facadeErrorMock) RefreshToken(userId uuid.UUID, refreshToken string) (*core.TokenResponseDTO, error) {
	refreshUserRecord := &refreshUserRecord{userId, refreshToken}
	facadeMock.refreshUserRecordArray = append(facadeMock.refreshUserRecordArray, refreshUserRecord)
	return nil, echo.ErrInternalServerError
}

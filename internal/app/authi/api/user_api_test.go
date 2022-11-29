package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

// CreateUserId Tests

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

// CreateUser Tests

func TestCreateUser_Successfully(t *testing.T) {
	facade := &facadeMock{createUserReturn: []error{nil}}
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
	// Exec
	err := userApi.CreateUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 1, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
	assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)

}

func TestCreateUser_CreateUser_ErrUserAlreadyExists(t *testing.T) {
	facade := &facadeMock{createUserReturn: []error{errormessages.ErrUserAlreadyExists}}
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
	// Exec
	err := userApi.CreateUser(c)
	// Assertions
	assert.Equal(t, echo.NewHTTPError(http.StatusConflict), err)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 1, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
	assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)

}

func TestCreateUser_CreateUser_InternalServerError(t *testing.T) {
	facade := &facadeMock{createUserReturn: []error{errSome}}
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
	// Exec
	err := userApi.CreateUser(c)
	// Assertions
	assert.Equal(t, echo.ErrInternalServerError, err)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 1, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.createUserRecordArray[0].userId)
	assert.Equal(t, password, facade.createUserRecordArray[0].authenticate.Password)
}

// LoginUser Tests

func TestLoginUser_Successfully(t *testing.T) {
	facade := &facadeMock{loginUserReturn: successfullyTokenResponse}
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
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.LoginUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 1, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.loginUserRecordArray[0].userId)
	assert.Equal(t, password, facade.loginUserRecordArray[0].authenticate.Password)
	assert.Equal(t, "{\"access_token\":\"some_access_token\",\"expires_in\":1,\"refresh_token\":\"some_refresh_token\",\"refresh_expires_in\":2}\n", rec.Body.String())

}

func TestLoginUser__LoginUser_ErrUnauthorized(t *testing.T) {
	facade := &facadeMock{loginUserReturn: errorTokenResponse}
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
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.LoginUser(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 1, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.loginUserRecordArray[0].userId)
	assert.Equal(t, password, facade.loginUserRecordArray[0].authenticate.Password)

}

// RefreshToken Tests

func TestRefreshToken_Successfully(t *testing.T) {
	facade := &facadeMock{refreshTokenReturn: successfullyTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, nil)
	req.Header.Set(authadapter.RefreshTokenHeaderName, refreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userRefreshPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.RefreshToken(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.refreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, facade.refreshTokenRecordArray[0].token)

}

func TestRefreshToken_RefreshToken_ErrUnauthorized(t *testing.T) {
	facade := &facadeMock{refreshTokenReturn: errorTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, userRootPath, nil)
	req.Header.Set(authadapter.RefreshTokenHeaderName, refreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam + userRefreshPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.RefreshToken(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 1, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.refreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, facade.refreshTokenRecordArray[0].token)
}

// Update Password Test

func TestUpdatePassword_Successfully(t *testing.T) {
	facade := &facadeMock{updatePasswordReturn: []error{nil}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPatch, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.UpdatePassword(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 1, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.updatePasswordRecordArray[0].userId)
	assert.Equal(t, password, facade.updatePasswordRecordArray[0].authenticate.Password)

}

func TestUpdatePassword_UpdatePassword_ErrUnauthorized(t *testing.T) {
	facade := &facadeMock{updatePasswordReturn: []error{errSome}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPatch, userRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.UpdatePassword(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 1, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.updatePasswordRecordArray[0].userId)
	assert.Equal(t, password, facade.updatePasswordRecordArray[0].authenticate.Password)
}

// Delete User Test

func TestDeleteUser_Successfully(t *testing.T) {
	facade := &facadeMock{deleteUserReturn: []error{nil}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, userRootPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.DeleteUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 1, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.deleteUserRecordArray[0].userId)
}

func TestDeleteUser_DeleteUser_ErrUnauthorized(t *testing.T) {
	facade := &facadeMock{deleteUserReturn: []error{errSome}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, userRootPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(userRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(authadapter.ClaimName, claimUser)
	// Exec
	err := userApi.DeleteUser(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.refreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.loginUserRecordArray))
	assert.Equal(t, 0, len(facade.createUserRecordArray))
	assert.Equal(t, 0, len(facade.updatePasswordRecordArray))
	assert.Equal(t, 1, len(facade.deleteUserRecordArray))
	assert.Equal(t, userId, facade.deleteUserRecordArray[0].userId)
}

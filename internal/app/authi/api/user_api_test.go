package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

// CreateUserId Tests

func TestCreateUserId_Successfully(t *testing.T) {
	userApi := &UserApi{&core.CoreMock{}}
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, adapter.AuthiRootPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	// Assertions
	if assert.NoError(t, userApi.CreateUserId(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// CreateUser Tests

func TestCreateUser_Successfully(t *testing.T) {
	facade := &core.CoreMock{CreateUserResponseArray: []*core.ErrorResponse{{Err: nil}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Exec
	err := userApi.CreateUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 1, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.CreateUserRecordArray[0].UserId)
	assert.Equal(t, password, facade.CreateUserRecordArray[0].Password)

}

func TestCreateUser_CreateUser_InternalServerError(t *testing.T) {
	facade := &core.CoreMock{CreateUserResponseArray: []*core.ErrorResponse{{Err: errSome}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	// Exec
	err := userApi.CreateUser(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 1, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.CreateUserRecordArray[0].UserId)
	assert.Equal(t, password, facade.CreateUserRecordArray[0].Password)
}

// LoginUser Tests

func TestLoginUser_Successfully(t *testing.T) {
	facade := &core.CoreMock{LoginUserResponseArray: successfullyTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam + adapter.AuthiLoginPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.LoginUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 1, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.LoginUserRecordArray[0].UserId)
	assert.Equal(t, password, facade.LoginUserRecordArray[0].Password)
	assert.Equal(t, "{\"access_token\":\"some_access_token\",\"expires_in\":1,\"refresh_token\":\"some_refresh_token\",\"refresh_expires_in\":2}\n", rec.Body.String())

}

func TestLoginUser__LoginUser_ErrUnauthorized(t *testing.T) {
	facade := &core.CoreMock{LoginUserResponseArray: errorTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam + adapter.AuthiLoginPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.LoginUser(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 1, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.LoginUserRecordArray[0].UserId)
	assert.Equal(t, password, facade.LoginUserRecordArray[0].Password)

}

// RefreshToken Tests

func TestRefreshToken_Successfully(t *testing.T) {
	facade := &core.CoreMock{RefreshTokenResponseArray: successfullyTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, nil)
	req.Header.Set(adapter.RefreshTokenHeaderName, refreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam + adapter.AuthiRefreshPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.RefreshToken(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.RefreshTokenRecordArray[0].UserId)
	assert.Equal(t, refreshToken, facade.RefreshTokenRecordArray[0].RefreshToken)

}

func TestRefreshToken_RefreshToken_ErrUnauthorized(t *testing.T) {
	facade := &core.CoreMock{RefreshTokenResponseArray: errorTokenResponse}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, adapter.AuthiRootPath, nil)
	req.Header.Set(adapter.RefreshTokenHeaderName, refreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam + adapter.AuthiRefreshPath)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.RefreshToken(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 1, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.RefreshTokenRecordArray[0].UserId)
	assert.Equal(t, refreshToken, facade.RefreshTokenRecordArray[0].RefreshToken)
}

// Update Password Test

func TestUpdatePassword_Successfully(t *testing.T) {
	facade := &core.CoreMock{UpdatePasswordResponseArray: []*core.ErrorResponse{{Err: nil}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPatch, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.UpdatePassword(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 1, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.UpdatePasswordRecordArray[0].UserId)
	assert.Equal(t, password, facade.UpdatePasswordRecordArray[0].Password)

}

func TestUpdatePassword_UpdatePassword_ErrUnauthorized(t *testing.T) {
	facade := &core.CoreMock{UpdatePasswordResponseArray: []*core.ErrorResponse{{Err: errSome}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPatch, adapter.AuthiRootPath, strings.NewReader(authenticationUserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.UpdatePassword(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 1, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.UpdatePasswordRecordArray[0].UserId)
	assert.Equal(t, password, facade.UpdatePasswordRecordArray[0].Password)
}

// Delete User Test

func TestDeleteUser_Successfully(t *testing.T) {
	facade := &core.CoreMock{DeleteUserResponseArray: []*core.ErrorResponse{{Err: nil}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, adapter.AuthiRootPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.DeleteUser(c)
	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 1, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.DeleteUserRecordArray[0].UserId)
}

func TestDeleteUser_DeleteUser_ErrUnauthorized(t *testing.T) {
	facade := &core.CoreMock{DeleteUserResponseArray: []*core.ErrorResponse{{Err: errSome}}}
	userApi := &UserApi{facade}
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodDelete, adapter.AuthiRootPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(loggerKey, log.WithField("Test", t.Name()))
	c.SetPath(adapter.AuthiRootPath + "/:" + userIdParam)
	c.SetParamNames(userIdParam)
	c.SetParamValues(userId.String())
	c.Set(adapter.ClaimName, claimUser)
	// Exec
	err := userApi.DeleteUser(c)
	// Assertions
	assert.Equal(t, echo.ErrUnauthorized, err)
	assert.Equal(t, 0, len(facade.RefreshTokenRecordArray))
	assert.Equal(t, 0, len(facade.LoginUserRecordArray))
	assert.Equal(t, 0, len(facade.CreateUserRecordArray))
	assert.Equal(t, 0, len(facade.UpdatePasswordRecordArray))
	assert.Equal(t, 1, len(facade.DeleteUserRecordArray))
	assert.Equal(t, userId, facade.DeleteUserRecordArray[0].UserId)
}

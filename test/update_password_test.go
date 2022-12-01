package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/test/util"
	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

const (
	someNewPassword = "NewPassword"
)

func TestUpdatePassword(t *testing.T) {
	token, userId := util.ObtainToken(t)
	status := util.UpdatePassword(userId, someNewPassword, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	_, status = util.Login(userId, someNewPassword)
	assert.Equal(t, status, http.StatusOK)
}

func TestUpdatePasswordWrongFormatAccessToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	status := util.UpdatePassword(userId, someNewPassword, "someWrongToken")
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestUpdatePasswordWrongUserIdToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	signKey := util.LoadPrivatKeyFile(util.PrivatKeyFile)
	customToken := util.CreateCustomJWTToken(uuid.NewString(), time.Now().Add(30*time.Minute).Unix(), signKey)
	status := util.UpdatePassword(userId, someNewPassword, customToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestUpdatePasswordWrongUserIdPath(t *testing.T) {
	token, _ := util.ObtainToken(t)
	status := util.UpdatePassword(uuid.NewString(), someNewPassword, token.AccessToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestUpdatePasswordExpiredToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	signKey := util.LoadPrivatKeyFile(util.PrivatKeyFile)
	customToken := util.CreateCustomJWTToken(userId, time.Now().Add(-1*time.Second).Unix(), signKey)
	status := util.UpdatePassword(userId, someNewPassword, customToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

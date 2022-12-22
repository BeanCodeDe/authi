package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/test/util"
	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

func TestDeleteUser_LoginUnable(t *testing.T) {
	token, userId := util.ObtainToken(t)
	status := util.DeleteUser(userId, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	_, status = util.Login(userId, util.DefaultPassword)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_TokenNotValid(t *testing.T) {
	token, userId := util.ObtainToken(t)
	status := util.DeleteUser(userId, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	_, status = util.RefreshToken(userId, token.AccessToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_WrongFormatAccessToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	status := util.DeleteUser(userId, "someWrongToken")
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_WrongUserIdToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	signKey := util.LoadPrivateKeyFile(util.PrivateKeyFile)
	customToken := util.CreateCustomJWTToken(uuid.NewString(), time.Now().Add(30*time.Minute).Unix(), signKey)
	status := util.DeleteUser(userId, customToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_WrongUserIdPath(t *testing.T) {
	token, _ := util.ObtainToken(t)
	status := util.DeleteUser(uuid.NewString(), token.AccessToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_ExpiredToken(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	signKey := util.LoadPrivateKeyFile(util.PrivateKeyFile)
	customToken := util.CreateCustomJWTToken(userId, time.Now().Add(-1*time.Second).Unix(), signKey)
	status := util.DeleteUser(userId, customToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestDeleteUser_Retry_LoginUnable(t *testing.T) {
	token, userId := util.ObtainToken(t)
	status := util.DeleteUser(userId, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	status = util.DeleteUser(userId, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	_, status = util.Login(userId, util.DefaultPassword)
	assert.Equal(t, status, http.StatusUnauthorized)
}

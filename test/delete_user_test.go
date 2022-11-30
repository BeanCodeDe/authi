package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
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

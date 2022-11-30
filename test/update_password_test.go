package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestUpdatePassword(t *testing.T) {
	token, userId := util.ObtainToken(t)
	someNewPassword := "NewPassword"
	status := util.UpdatePassword(userId, someNewPassword, token.AccessToken)
	assert.Equal(t, status, http.StatusNoContent)
	_, status = util.Login(userId, someNewPassword)
	assert.Equal(t, status, http.StatusOK)
}

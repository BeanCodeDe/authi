package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestLogin(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	_, status := util.Login(userId, util.DefaultPassword)
	assert.Equal(t, status, http.StatusOK)
}

func TestLoginFailed(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	_, status := util.Login(userId, "wrongPassword")
	assert.Equal(t, status, http.StatusUnauthorized)
}

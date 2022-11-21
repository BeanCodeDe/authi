package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestLogin(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	authenticate := &util.Authenticate{Password: util.DefaultPassword}
	_, status := util.Login(userId, authenticate)
	assert.Equal(t, status, http.StatusOK)
}

func TestLoginFailed(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	authenticate := &util.Authenticate{Password: "wrongPassword"}
	_, status := util.Login(userId, authenticate)
	assert.Equal(t, status, http.StatusUnauthorized)
}

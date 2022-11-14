package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestLogin(t *testing.T) {
	userId := util.CreateUserForFurtherTesting(t)
	user := &util.UserDTO{ID: userId, Password: util.DefaultPassword}
	_, status := util.Login(user)
	assert.Equal(t, status, http.StatusOK)
}

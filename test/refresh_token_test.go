package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestAuth(t *testing.T) {
	token := util.OptainToken(t)
	newToken, status := util.RefreshToken(token.AccessToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEqual(t, newToken, nil)
}

func TestAuthWrongAccessToken(t *testing.T) {
	token := util.OptainToken(t)
	_, status := util.RefreshToken(token.RefreshToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestAuthWrongRefreshToken(t *testing.T) {
	token := util.OptainToken(t)
	_, status := util.RefreshToken(token.AccessToken, token.AccessToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

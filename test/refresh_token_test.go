package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/test/util"
	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

func TestAuth(t *testing.T) {
	token, _ := util.OptainToken(t)
	newToken, status := util.RefreshToken(token.AccessToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEqual(t, newToken, nil)
}

func TestAuthWrongFormatAccessToken(t *testing.T) {
	token, _ := util.OptainToken(t)
	_, status := util.RefreshToken(token.RefreshToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestAuthWrongUserIdToken(t *testing.T) {
	token, _ := util.OptainToken(t)
	signKey := util.LoadPrivatKeyFile(util.PrivatKeyFile)
	customToken := util.CreateCustomJWTToken(uuid.NewString(), time.Now().Add(30*time.Minute).Unix(), signKey)
	_, status := util.RefreshToken(customToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestAuthExpiredToken(t *testing.T) {
	token, user := util.OptainToken(t)
	signKey := util.LoadPrivatKeyFile(util.PrivatKeyFile)
	customToken := util.CreateCustomJWTToken(user.ID, time.Now().Add(-1*time.Second).Unix(), signKey)
	_, status := util.RefreshToken(customToken, token.RefreshToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

func TestAuthWrongRefreshToken(t *testing.T) {
	token, _ := util.OptainToken(t)
	_, status := util.RefreshToken(token.AccessToken, token.AccessToken)
	assert.Equal(t, status, http.StatusUnauthorized)
}

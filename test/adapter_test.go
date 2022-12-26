package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/BeanCodeDe/authi/test/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdapterRefreshToken(t *testing.T) {
	t.Setenv(adapter.EnvAuthUrl, util.Url)
	token, userId := util.ObtainToken(t)
	authi := adapter.NewAuthiAdapter(uuid.NewString())
	refreshedToken, err := authi.RefreshToken(userId, token.AccessToken, token.RefreshToken)

	assert.Nil(t, err)
	assert.NotNil(t, refreshedToken)

	status := util.DeleteUser(userId, refreshedToken.AccessToken)
	assert.Equal(t, http.StatusNoContent, status)
}

func TestAdapterLogin(t *testing.T) {
	t.Setenv(adapter.EnvAuthUrl, util.Url)
	userId := util.CreateUserForFurtherTesting(t)
	authi := adapter.NewAuthiAdapter(uuid.NewString())
	refreshedToken, err := authi.GetToken(userId, util.DefaultPassword)

	assert.Nil(t, err)
	assert.NotNil(t, refreshedToken)

	status := util.DeleteUser(userId, refreshedToken.AccessToken)
	assert.Equal(t, http.StatusNoContent, status)
}

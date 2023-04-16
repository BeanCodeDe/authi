package core

import (
	"errors"
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	userId         = uuid.New()
	refreshToken   = "someRefreshToken"
	password       = "some password"
	authenticate   = &adapter.AuthenticateDTO{Password: password}
	privateKeyPath = "../../../../deployments/data/token/jwtRS256.key"
	errUnknown     = errors.New("some error")
)

func TestRandomString(t *testing.T) {
	randomString := randomString()
	assert.Len(t, randomString, 32)
	assert.Regexp(t, "[a-zA-Z0-9]*", randomString)
}

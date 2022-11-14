package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"github.com/google/uuid"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateId(t *testing.T) {
	userId, status := util.CreateUserId()
	assert.Equal(t, status, http.StatusOK)
	_, err := uuid.Parse(userId)
	assert.Equal(t, err, nil)
}

func TestCreateUser(t *testing.T) {
	userId, status := util.CreateUserId()
	assert.Equal(t, status, http.StatusOK)
	_, err := uuid.Parse(userId)
	assert.Equal(t, err, nil)

	status = util.CreateUser(userId, "random_password")
	assert.Equal(t, status, http.StatusCreated)
}

func TestCreateUserRetry(t *testing.T) {
	userId, status := util.CreateUserId()
	assert.Equal(t, status, http.StatusOK)
	_, err := uuid.Parse(userId)
	assert.Equal(t, err, nil)

	status = util.CreateUser(userId, "random_password")
	assert.Equal(t, status, http.StatusCreated)

	status = util.CreateUser(userId, "random_password")
	assert.Equal(t, status, http.StatusConflict)
}

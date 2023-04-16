package test

import (
	"net/http"
	"testing"

	"github.com/BeanCodeDe/authi/test/util"
	"gopkg.in/go-playground/assert.v1"
)

func TestInitUserLogin(t *testing.T) {
	_, status := util.Login("c5ffc340-507e-4c66-a6ce-a7d98842f9ba", "someSecretPassword")
	assert.Equal(t, status, http.StatusOK)
	_, secondStatus := util.Login("5cc3621d-e5ac-4d81-93df-462b27e0cc2b", "someOtherPassword")
	assert.Equal(t, secondStatus, http.StatusOK)
}

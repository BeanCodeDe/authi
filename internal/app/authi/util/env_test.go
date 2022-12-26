package util

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvWithFallback_Successfully(t *testing.T) {
	someEnv := "SOME_ENV"
	someEnvValue := "SOME VALUE"
	someFallbackValue := "SOME VALUE"
	t.Setenv(someEnv, someEnvValue)

	value := GetEnvWithFallback(someEnv, someFallbackValue)
	assert.Equal(t, someEnvValue, value)
}

func TestGetEnvWithFallback_NotFound(t *testing.T) {
	someEnv := "SOME_ENV"
	someFallbackValue := "SOME VALUE"

	value := GetEnvWithFallback(someEnv, someFallbackValue)
	assert.Equal(t, someFallbackValue, value)
}

func TestGetEnv_Successfully(t *testing.T) {
	someEnv := "SOME_ENV"
	someEnvValue := "SOME VALUE"
	t.Setenv(someEnv, someEnvValue)

	value, err := GetEnv(someEnv)
	assert.Nil(t, err)
	assert.Equal(t, someEnvValue, value)
}

func TestGetEnv_NotFound(t *testing.T) {
	someEnv := "SOME_ENV"

	value, err := GetEnv(someEnv)
	assert.Empty(t, value)
	assert.ErrorContains(t, err, "not found")
}

func TestGetEnvIntWithFallback_Successfully(t *testing.T) {
	someEnv := "SOME_ENV"
	someEnvValue := 5
	someFallbackValue := 10
	t.Setenv(someEnv, strconv.Itoa(someEnvValue))

	value, err := GetEnvIntWithFallback(someEnv, someFallbackValue)
	assert.Nil(t, err)
	assert.Equal(t, someEnvValue, value)
}

func TestGetEnvIntWithFallback_NotFound(t *testing.T) {
	someEnv := "SOME_ENV"
	someFallbackValue := 10

	value, err := GetEnvIntWithFallback(someEnv, someFallbackValue)
	assert.Nil(t, err)
	assert.Equal(t, someFallbackValue, value)
}

func TestGetEnvIntWithFallback_WrongFormat(t *testing.T) {
	someEnv := "SOME_ENV"
	someEnvValue := "five"
	someFallbackValue := 10
	t.Setenv(someEnv, someEnvValue)

	value, err := GetEnvIntWithFallback(someEnv, someFallbackValue)
	assert.Equal(t, 0, value)
	assert.ErrorContains(t, err, "invalid syntax")
}

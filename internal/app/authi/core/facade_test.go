package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	randomString := randomString()
	assert.Len(t, randomString, 32)
	assert.Regexp(t, "[a-zA-Z0-9]*", randomString)
}

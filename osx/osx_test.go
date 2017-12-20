package osx

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("no value", GetEnvOr("TestGetEnv", "no value"))
	assert.NoError(os.Setenv("TestGetEnv", "value"))
	assert.Equal("value", GetEnvOr("TestGetEnv", "no value"))
}


package osx

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	assert := require.New(t)

	assert.Equal("no value", GetEnvOr("TestGetEnv", "no value"))
	assert.NoError(os.Setenv("TestGetEnv", "value"))
	assert.Equal("value", GetEnvOr("TestGetEnv", "no value"))
}

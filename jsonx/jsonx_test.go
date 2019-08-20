package jsonx_test

import (
	"github.com/chen56/go-common/jsonx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMustMarshal2String(t *testing.T) {
	require.Equal(t, `"s"`, jsonx.MustMarshal2String("s"))
}

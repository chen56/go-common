package must

import (
	"github.com/chen56/go-common/reflectx"
	"github.com/stretchr/testify/require"
	"testing"
)

type X struct {
}

func TestIsEmpty(t *testing.T) {
	var x X
	require.True(t, reflectx.IsZero(x))
	require.True(t, reflectx.IsZero(&x))
}

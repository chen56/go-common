package must_test

import (
	"github.com/chen56/go-common/must"
	"testing"
)

func TestNil(t *testing.T) {
	var x *testing.T = nil
	must.Nil(x)
}

package reflectx_test

import (
	"fmt"
	"github.com/chen56/go-common/reflectx"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

type X struct {
	time time.Time
	Y    Y
}
type Y struct {
	str string
}

func TestIsEmpty(t *testing.T) {
	var x X
	require.True(t, reflectx.IsZero(x))
	require.True(t, reflectx.IsZero(x.time))
}
func TestVisitFields(t *testing.T) {
	err := reflectx.VisitFields(X{
		time: time.Now(),
		Y: Y{
			str: "sss",
		},
	}, func(field reflect.Type) error {
		fmt.Println("", field.String())
		return nil
	})
	require.NoError(t, err)
}

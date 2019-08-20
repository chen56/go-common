package timex_test

import (
	"fmt"
	"github.com/chen56/go-common/timex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLocation(t *testing.T) {
	a := require.New(t)

	x := time.Date(1980, 1, 1, 0, 0, 0, 0, timex.LocationAsiaShanghai)

	a.Equal("1980-01-01T00:00:00+08:00", x.Format(time.RFC3339))

	fmt.Println(timex.Zero)
}

func TestMustParse(t *testing.T) {
	a := require.New(t)

	x := timex.MustParse(time.RFC3339, "1980-01-01T00:00:00+08:00")
	a.Equal("1980-01-01T00:00:00+08:00", x.Format(time.RFC3339))
}

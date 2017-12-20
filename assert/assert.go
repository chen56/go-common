// Deprecated: replaced by package: must
package assert

import (
	"fmt"
)

func True(b bool) {
	if !b {
		panic(fmt.Sprintf("expected true"))
	}
}

func False(b bool) {
	if b {
		panic(fmt.Sprintf("expected false"))
	}
}

func Falsef(b bool, msg string, args ...interface{}) {
	if b {
		panic(fmt.Sprintf(msg, args...))
	}
}

func Fail(msg string) {
	panic(msg)
}

func Failf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

func NoErr(err error) {
	if err != nil {
		// Handle error
		panic(err)
	}
}

func NotNil(x interface{}) {
	if x == nil {
		panic(x)
	}
}

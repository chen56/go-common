package runtimex

import (
	"flag"
)

func IsRunInTest() bool {
	return flag.Lookup("test.v") != nil
}

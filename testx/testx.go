package testx

import (
	"runtime"
	"strings"
	"fmt"
	"testing"
)

func CurrentFuncName()string {
	return traceFuncName(3)
}

func traceFuncName(skip int)string{
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(skip, pc)
	f := runtime.FuncForPC(pc[0])
	x:=strings.Split(f.Name(),".")
	return x[len(x)-1]
}

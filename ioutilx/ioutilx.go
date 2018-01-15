package ioutilx

import (
	"io/ioutil"

	"github.com/chen56/go-common/must"
)

func MustReadFileToString(file string) string {
	content, err := ioutil.ReadFile(file)
	must.NoErr(err)
	return string(content)
}

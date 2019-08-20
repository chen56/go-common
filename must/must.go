package must

import (
	"github.com/chen56/go-common/reflectx"
	"github.com/pkg/errors"
)

func NoError(err error, msgAndArgs ...interface{}) {
	if err != nil {
		panic(errors.Wrapf(err, messageFromMsgAndArgs("must NoError", msgAndArgs...)))
	}
}

func NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	if reflectx.IsZero(object) {
		panic(messageFromMsgAndArgs("must NotEmpty", msgAndArgs...))
	}
}
func NotNil(object interface{}, msgAndArgs ...interface{}) {
	if isNil(object) {
		panic(messageFromMsgAndArgs("must NotNil", msgAndArgs...))
	}
}
func Nil(object interface{}, msgAndArgs ...interface{}) {
	if !isNil(nil) {
		panic(messageFromMsgAndArgs("must Nil", msgAndArgs...))
	}
}
func True(b bool, msgAndArgs ...interface{}) {
	if !b {
		panic(messageFromMsgAndArgs("must true", msgAndArgs...))
	}
}
func False(b bool, msgAndArgs ...interface{}) {
	if b {
		panic(messageFromMsgAndArgs("must false", msgAndArgs...))
	}
}
func Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if !objectsAreEqual(expected, actual) {
		panic(messageFromMsgAndArgs("must equals", msgAndArgs...))
	}
}
func NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if objectsAreEqual(expected, actual) {
		panic(messageFromMsgAndArgs("must NotEquals", msgAndArgs...))
	}
}

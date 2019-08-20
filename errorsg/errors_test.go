package errorsg_test

import (
	"context"
	"fmt"
	"github.com/chen56/go-common/errorsg"
	"github.com/c2fo/testify/require"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Example_it_is_error() {
	var err error = errorsg.NotFound("user不存在: %s", "chen").Err()
	fmt.Println(err)

	// Output:
	// NotFound: user不存在: chen
}
func Example_retryInfo() {
	fmt.Println(toJson(errorsg.Unavailable("暂时不能服务，稍后再试").WithRetryInfo(&errdetails.RetryInfo{
		RetryDelay: &duration.Duration{
			Seconds: 1,
		},
	})))
	// Output:
	// {"error":"暂时不能服务，稍后再试","message":"暂时不能服务，稍后再试","code":14,"details":[{"@type":"type.googleapis.com/google.rpc.RetryInfo","retryDelay":"1s"}]}
}
func TestLog_panic_error(t *testing.T) {
	_ = errorsg.Try(func() error {
		panic(errors.WithStack(errors.New("error")))
	})
}
func TestLog_panic_string(t *testing.T) {
	_ = errorsg.Try(func() error {
		panic("string Error")
	})
}
func TestLog_error(t *testing.T) {
	_ = errorsg.Try(func() error {
		return errors.New("error")
	})
}
func TestLog_errorWithStack(t *testing.T) {
	_ = errorsg.Try(func() error {
		return errors.WithStack(errors.New("error with stack"))
	})
}

func TestLog_statusError(t *testing.T) {
	_ = errorsg.Try(func() error {
		return errorsg.Unauthenticated("Unauthenticated").Err()
	})
}
func TestLog_statusError_maybeSystemfailure(t *testing.T) {
	_ = errorsg.Try(func() error {
		return errorsg.Unavailable("xxx Unavailable").Err()
	})
}

func TestRetryInfo(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		err      *errorsg.Error
		duration *duration.Duration
	}{
		{
			err:      errorsg.Internal("x"),
			duration: &duration.Duration{Seconds: 1},
		},
		{
			err:      errorsg.Unavailable("x"),
			duration: &duration.Duration{Seconds: 1},
		},
		{
			err:      errorsg.ResourceExhausted("x"),
			duration: &duration.Duration{Seconds: 60},
		},
		{
			err:      errorsg.Unknown("x"),
			duration: nil,
		},
		{
			err:      errorsg.Unauthenticated("x"),
			duration: nil,
		},
	}

	for i, test := range tests {
		err := errorsg.Try(func() error {
			return test.err.Err()
		})
		if test.duration == nil {
			require.Nil(t, errorsg.FromError(err).RetryInfo(), "%d - %s retry Nil", i, test.err)
		} else {
			require.Equal(t, test.duration.Seconds, errorsg.FromError(err).RetryInfo().RetryDelay.Seconds, "%d - %s retry Equal", i, test.err)
		}
	}
}
func TestTry(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		msg      string
		f        func() error
		expected *errorsg.Error
	}{
		{
			msg: "return : status",
			f: func() error {
				return errorsg.Internal("x").Err()
			},
			expected: errorsg.Internal("x"),
		},
		{
			msg: "return : nil",
			f: func() error {
				return nil
			},
			expected: nil,
		},
		{
			msg: "return : error",
			f: func() error {
				return errors.New("json error")
			},
			expected: errorsg.Unknown("catch err"),
		},
		{
			msg: "recover : string",
			f: func() error {
				panic("json error")
			},
			expected: errorsg.Unknown("catch recover err"),
		},
		{
			msg: "recover : error",
			f: func() error {
				panic(errors.New("json error"))
			},
			expected: errorsg.Unknown("catch recover err"),
		},
		{
			msg: "recover : status",
			f: func() error {
				panic(errorsg.FailedPrecondition("audio invalid"))
			},
			expected: errorsg.FailedPrecondition("audio invalid"),
		},
	}

	for i, test := range tests {
		err := errorsg.Try(test.f)
		if test.expected == nil {
			require.Nil(t, err, "%d - %s - Nil", i, test.msg)
		} else {
			require.Equal(t, test.expected.Err().Error(), err.Error(), "%d - %s - Equal", i, test.msg)
			require.Equal(t, test.expected.Code(), errorsg.FromError(err).Code(), "%d - %s - Equal", i, test.msg)
		}
	}
}

func toJson(err *errorsg.Error) string {
	ctx := context.Background()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("", "", nil) // Pass in an empty request to match the signature
	runtime.DefaultHTTPError(ctx, &runtime.ServeMux{}, &runtime.JSONPb{}, w, req, err.Err())
	return w.Body.String()
}

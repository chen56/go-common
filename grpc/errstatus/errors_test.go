package errstatus_test

import (
	"context"
	"fmt"
	"github.com/c2fo/testify/require"
	"github.com/chen56/go-common/grpc/errstatus"
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
	var err error = errstatus.NotFound("user不存在: %s", "chen").Err()
	fmt.Println(err)

	// Output:
	// NotFound: user不存在: chen
}
func Example_retryInfo() {
	fmt.Println(toJson(errstatus.Unavailable("暂时不能服务，稍后再试").WithRetryInfo(&errdetails.RetryInfo{
		RetryDelay: &duration.Duration{
			Seconds: 1,
		},
	})))
	// Output:
	// {"error":"暂时不能服务，稍后再试","message":"暂时不能服务，稍后再试","code":14,"details":[{"@type":"type.googleapis.com/google.rpc.RetryInfo","retryDelay":"1s"}]}
}
func TestLog_panic_error(t *testing.T) {
	_ = errstatus.Try(func() error {
		panic(errors.WithStack(errors.New("error")))
	})
}
func TestLog_panic_string(t *testing.T) {
	_ = errstatus.Try(func() error {
		panic("string ErrorStatus")
	})
}
func TestLog_error(t *testing.T) {
	_ = errstatus.Try(func() error {
		return errors.New("error")
	})
}
func TestLog_errorWithStack(t *testing.T) {
	_ = errstatus.Try(func() error {
		return errors.WithStack(errors.New("error with stack"))
	})
}

func TestLog_statusError(t *testing.T) {
	_ = errstatus.Try(func() error {
		return errstatus.Unauthenticated("Unauthenticated").Err()
	})
}
func TestLog_statusError_maybeSystemfailure(t *testing.T) {
	_ = errstatus.Try(func() error {
		return errstatus.Unavailable("xxx Unavailable").Err()
	})
}

func TestRetryInfo(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		err      *errstatus.ErrorStatus
		duration *duration.Duration
	}{
		{
			err:      errstatus.Internal("x"),
			duration: &duration.Duration{Seconds: 1},
		},
		{
			err:      errstatus.Unavailable("x"),
			duration: &duration.Duration{Seconds: 1},
		},
		{
			err:      errstatus.ResourceExhausted("x"),
			duration: &duration.Duration{Seconds: 60},
		},
		{
			err:      errstatus.Unknown("x"),
			duration: nil,
		},
		{
			err:      errstatus.Unauthenticated("x"),
			duration: nil,
		},
	}

	for i, test := range tests {
		err := errstatus.Try(func() error {
			return test.err.Err()
		})
		if test.duration == nil {
			require.Nil(t, errstatus.FromError(err).RetryInfo(), "%d - %s retry Nil", i, test.err)
		} else {
			require.Equal(t, test.duration.Seconds, errstatus.FromError(err).RetryInfo().RetryDelay.Seconds, "%d - %s retry Equal", i, test.err)
		}
	}
}
func TestTry(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		msg      string
		f        func() error
		expected *errstatus.ErrorStatus
	}{
		{
			msg: "return : status",
			f: func() error {
				return errstatus.Internal("x").Err()
			},
			expected: errstatus.Internal("x"),
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
			expected: errstatus.Unknown("catch err"),
		},
		{
			msg: "recover : string",
			f: func() error {
				panic("json error")
			},
			expected: errstatus.Unknown("catch recover err"),
		},
		{
			msg: "recover : error",
			f: func() error {
				panic(errors.New("json error"))
			},
			expected: errstatus.Unknown("catch recover err"),
		},
		{
			msg: "recover : status",
			f: func() error {
				panic(errstatus.FailedPrecondition("audio invalid"))
			},
			expected: errstatus.FailedPrecondition("audio invalid"),
		},
	}

	for i, test := range tests {
		err := errstatus.Try(test.f)
		if err == nil {
			require.Nil(t, test.expected, "%d - %s - Nil", i, test.msg)
		} else {
			require.Equal(t, test.expected.Err().Error(), err.Error(), "%d - %s - Equal", i, test.msg)
			require.Equal(t, test.expected.Code(), errstatus.FromError(err).Code(), "%d - %s - Equal", i, test.msg)
		}
	}
}

func toJson(err *errstatus.ErrorStatus) string {
	ctx := context.Background()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("", "", nil) // Pass in an empty request to match the signature
	runtime.DefaultHTTPError(ctx, &runtime.ServeMux{}, &runtime.JSONPb{}, w, req, err.Err())
	return w.Body.String()
}

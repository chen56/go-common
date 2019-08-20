package statusx

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"runtime/debug"
	"time"
)

func (x *Error) setRetryIfNeed() *Error {
	for _, d := range x.GRPCStatus().Details() {
		// 如果已经设置retryInfo，就不设了
		if _, ok := d.(*errdetails.RetryInfo); ok {
			return x
		}
	}

	code := x.status.Code()
	switch code {
	case codes.Aborted, codes.Internal, codes.Unavailable:
		return x.WithRetryInfo(&errdetails.RetryInfo{
			RetryDelay: ptypes.DurationProto(time.Second),
		})
	case codes.ResourceExhausted:
		return x.WithRetryInfo(&errdetails.RetryInfo{
			RetryDelay: ptypes.DurationProto(60 * time.Second),
		})
	default:
		return x
	}
}

func Try(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fromErr("catch recover err", r)
			return
		}
	}()
	err = fromErr("catch err", f())
	return
}

func fromErr(title string, e interface{}) error {
	if e != nil {
		// grpcStatus 认为是已处理过的错误
		if e, ok := e.(grpcStatus); ok {
			code := e.GRPCStatus().Code()
			if code == codes.Unknown || code == codes.Internal || code == codes.Unavailable {
				d := e.GRPCStatus().Details()
				stack := string(debug.Stack())
				logrus.Errorf("%s: code: %s : %s\ndetails: %+v\nstack: %s", title, code.String(), e, d, stack)
				fmt.Printf("%s: code: %s : %s\ndetails: %+v\nstack: %s", title, code.String(), e, d, stack)
			} else if logrus.IsLevelEnabled(logrus.DebugLevel) {
				d := e.GRPCStatus().Details()
				stack := string(debug.Stack())
				logrus.Debugf("%s: code: %s : %s\ndetails: %+v\nstack: %s", title, code.String(), e, d, stack)
				fmt.Printf("%s: code: %s : %s\ndetails: %+v\nstack: %s", title, code.String(), e, d, stack)
			}
			return FromStatus(e.GRPCStatus()).
				setRetryIfNeed().Err()
		}
		stack := causeStackTrace(e)
		logrus.Errorf("%s: %s %+v", title, e, stack)
		fmt.Printf("%s: %s %+v\n", title, e, stack)
		return Unknown(title).
			WithDebugMessage("%s: %s\n%+v", title, e, stack).
			setRetryIfNeed().Err()
	}
	return nil
}

// 查找err的cause链，找到最后一个有StackTrace的深层错误的StackTrace返回之
func causeStackTrace(err interface{}) string {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	type causer interface {
		Cause() error
	}
	if err, ok := err.(error); ok {
		var errs = []error{err}
		for err != nil {
			cause, ok := err.(causer)
			if !ok {
				break
			}
			err = cause.Cause()
			errs = append(errs, err)
		}
		for i := len(errs) - 1; i >= 0; i-- {
			hasStack, ok := errs[i].(stackTracer)
			if ok {
				return fmt.Sprintf("%+v", hasStack.StackTrace())
			}
		}
	}
	return string(debug.Stack())
}

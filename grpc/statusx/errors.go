// grpc 的错误处理包
// 包含：
//   - 错误生成的帮助方法 : return errstatus.PermissionDenied("token error")
//   - grpc拦截器 grpc.NewUnaryServerInterceptor

package statusx

import (
	"fmt"
	"github.com/chen56/go-common/must"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcStatus interface {
	GRPCStatus() *status.Status
}

func FromStatus(status *status.Status) *Status {
	must.True(status.Code() != codes.OK)
	return &Status{
		status: *status,
	}
}

func FromError(err error) *Status {
	must.NotNil(err)
	return MapIfStdError(err, func(err error) *Status {
		return Unknown("Unknown error").WithDebugError(err)
	})
}

// MapIfStdError 如果参数err是普通的error(而不是errorsg包的Error),则 用newError转换为新error
func MapIfStdError(err error, newError func(err error) *Status) *Status {
	must.NotNil(err)
	if e, ok := err.(*statusError); ok {
		return e.errorStatus
	}
	if e, ok := err.(grpcStatus); ok {
		return FromStatus(e.GRPCStatus())
	}
	return newError(err)
}

func Canceled(msgAndArgs ...interface{}) *Status {
	return newError(codes.Canceled, msgAndArgs)
}
func Unknown(msgAndArgs ...interface{}) *Status {
	return newError(codes.Unknown, msgAndArgs)
}
func InvalidArgument(msgAndArgs ...interface{}) *Status {
	return newError(codes.InvalidArgument, msgAndArgs)
}
func DeadlineExceeded(msgAndArgs ...interface{}) *Status {
	return newError(codes.DeadlineExceeded, msgAndArgs)
}
func NotFound(msgAndArgs ...interface{}) *Status {
	return newError(codes.NotFound, msgAndArgs)
}
func AlreadyExists(msgAndArgs ...interface{}) *Status {
	return newError(codes.AlreadyExists, msgAndArgs)
}
func PermissionDenied(msgAndArgs ...interface{}) *Status {
	return newError(codes.PermissionDenied, msgAndArgs)
}
func ResourceExhausted(msgAndArgs ...interface{}) *Status {
	return newError(codes.ResourceExhausted, msgAndArgs)
}
func FailedPrecondition(msgAndArgs ...interface{}) *Status {
	return newError(codes.FailedPrecondition, msgAndArgs)
}
func Aborted(msgAndArgs ...interface{}) *Status {
	return newError(codes.Aborted, msgAndArgs)
}
func OutOfRange(msgAndArgs ...interface{}) *Status {
	return newError(codes.OutOfRange, msgAndArgs)
}
func Unimplemented(msgAndArgs ...interface{}) *Status {
	return newError(codes.Unimplemented, msgAndArgs)
}
func Internal(msgAndArgs ...interface{}) *Status {
	return newError(codes.Internal, msgAndArgs)
}
func Unavailable(msgAndArgs ...interface{}) *Status {
	return newError(codes.Unavailable, msgAndArgs)
}
func DataLoss(msgAndArgs ...interface{}) *Status {
	return newError(codes.DataLoss, msgAndArgs)
}
func Unauthenticated(msgAndArgs ...interface{}) *Status {
	return newError(codes.Unauthenticated, msgAndArgs)
}

type Status struct {
	status status.Status
}

// to error interface
func (x *Status) Err() error {
	return &statusError{errorStatus: x}
}

func (x *Status) Code() codes.Code {
	return x.status.Code()
}
func (x *Status) Message() string {
	return x.status.Message()
}
func (x *Status) GRPCStatus() *status.Status {
	return &x.status
}

func (x *Status) RetryInfo() *errdetails.RetryInfo {
	for _, d := range x.GRPCStatus().Details() {
		if retry, ok := d.(*errdetails.RetryInfo); ok {
			return retry
		}
	}
	return nil
}

func (x *Status) WithDebugMessage(debugMsgAndArgs ...interface{}) *Status {
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: format("[debug]", debugMsgAndArgs),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Status) WithDetails(details ...proto.Message) *Status {
	newStatus, err := x.status.WithDetails(details...)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Status) WithRetryInfo(retryInfo *errdetails.RetryInfo) *Status {
	newStatus, err := x.status.WithDetails(retryInfo)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Status) WithDebugError(err error) *Status {
	must.NotNil(err)
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: err.Error(),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func newError(code codes.Code, msgAndArgs []interface{}) *Status {
	return &Status{
		status: *status.New(code, format(code.String(), msgAndArgs)),
	}
}

func format(defaultMsg string, msgAndArgs []interface{}) string {
	if len(msgAndArgs) == 0 {
		return defaultMsg
	}

	msg, ok := msgAndArgs[0].(string)
	if !ok {
		return fmt.Sprint(msgAndArgs...)
	}

	if len(msgAndArgs) == 1 {
		return msg
	}
	return fmt.Sprintf(msg, msgAndArgs[1:]...)
}

// grpc 的错误处理包
// 包含：
//   - 错误生成的帮助方法 : return errorsg.PermissionDenied("token error")
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

func FromStatus(status *status.Status) *Error {
	must.True(status.Code() != codes.OK)
	return &Error{
		status: *status,
	}
}

func FromError(err error) *Error {
	return MapIfStdError(err, func(err error) *Error {
		return Unknown("Unknown error").WithDebugError(err)
	})
}

// MapIfStdError 如果参数err是普通的error(而不是errorsg包的Error),则 用newError转换为新error
func MapIfStdError(err error, newError func(err error) *Error) *Error {
	must.NotNil(err)
	if e, ok := err.(*statusError); ok {
		return e.errorStatus
	}
	if e, ok := err.(grpcStatus); ok {
		return FromStatus(e.GRPCStatus())
	}
	return newError(err)
}

func Canceled(msgAndArgs ...interface{}) *Error {
	return newError(codes.Canceled, msgAndArgs)
}
func Unknown(msgAndArgs ...interface{}) *Error {
	return newError(codes.Unknown, msgAndArgs)
}
func InvalidArgument(msgAndArgs ...interface{}) *Error {
	return newError(codes.InvalidArgument, msgAndArgs)
}
func DeadlineExceeded(msgAndArgs ...interface{}) *Error {
	return newError(codes.DeadlineExceeded, msgAndArgs)
}
func NotFound(msgAndArgs ...interface{}) *Error {
	return newError(codes.NotFound, msgAndArgs)
}
func AlreadyExists(msgAndArgs ...interface{}) *Error {
	return newError(codes.AlreadyExists, msgAndArgs)
}
func PermissionDenied(msgAndArgs ...interface{}) *Error {
	return newError(codes.PermissionDenied, msgAndArgs)
}
func ResourceExhausted(msgAndArgs ...interface{}) *Error {
	return newError(codes.ResourceExhausted, msgAndArgs)
}
func FailedPrecondition(msgAndArgs ...interface{}) *Error {
	return newError(codes.FailedPrecondition, msgAndArgs)
}
func Aborted(msgAndArgs ...interface{}) *Error {
	return newError(codes.Aborted, msgAndArgs)
}
func OutOfRange(msgAndArgs ...interface{}) *Error {
	return newError(codes.OutOfRange, msgAndArgs)
}
func Unimplemented(msgAndArgs ...interface{}) *Error {
	return newError(codes.Unimplemented, msgAndArgs)
}
func Internal(msgAndArgs ...interface{}) *Error {
	return newError(codes.Internal, msgAndArgs)
}
func Unavailable(msgAndArgs ...interface{}) *Error {
	return newError(codes.Unavailable, msgAndArgs)
}
func DataLoss(msgAndArgs ...interface{}) *Error {
	return newError(codes.DataLoss, msgAndArgs)
}
func Unauthenticated(msgAndArgs ...interface{}) *Error {
	return newError(codes.Unauthenticated, msgAndArgs)
}

type Error struct {
	status status.Status
}

type statusError struct {
	errorStatus *Error
}

// impl error interface
func (x *statusError) Error() string {
	return fmt.Sprintf("%s: %s", x.errorStatus.Code(), x.errorStatus.Message())
}

func (x *statusError) GRPCStatus() *status.Status {
	return &x.errorStatus.status
}

// to error interface
func (x *Error) Err() error {
	return &statusError{errorStatus: x}
}

func (x *Error) Code() codes.Code {
	return x.status.Code()
}
func (x *Error) Message() string {
	return x.status.Message()
}
func (x *Error) GRPCStatus() *status.Status {
	return &x.status
}

func (x *Error) RetryInfo() *errdetails.RetryInfo {
	for _, d := range x.GRPCStatus().Details() {
		if retry, ok := d.(*errdetails.RetryInfo); ok {
			return retry
		}
	}
	return nil
}

func (x *Error) WithDebugMessage(debugMsgAndArgs ...interface{}) *Error {
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: format("[debug]", debugMsgAndArgs),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Error) WithDetails(details ...proto.Message) *Error {
	newStatus, err := x.status.WithDetails(details...)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Error) WithRetryInfo(retryInfo *errdetails.RetryInfo) *Error {
	newStatus, err := x.status.WithDetails(retryInfo)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *Error) WithDebugError(err error) *Error {
	must.NotNil(err)
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: err.Error(),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func newError(code codes.Code, msgAndArgs []interface{}) *Error {
	return &Error{
		status: *status.New(code, format(code.String(), msgAndArgs)),
	}
}

func format(defaultMsg string, msgAndArgs []interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
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

// grpc 的错误处理包
// 包含：
//   - 错误生成的帮助方法 : return errstatus.PermissionDenied("token error")
//   - grpc拦截器 grpc.NewUnaryServerInterceptor

package errstatus

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

func FromStatus(status *status.Status) *ErrorStatus {
	must.True(status.Code() != codes.OK)
	return &ErrorStatus{
		status: *status,
	}
}

func FromError(err error) *ErrorStatus {
	return MapIfStdError(err, func(err error) *ErrorStatus {
		return Unknown("Unknown error").WithDebugError(err)
	})
}

// MapIfStdError 如果参数err是普通的error(而不是errorsg包的Error),则 用newError转换为新error
func MapIfStdError(err error, newError func(err error) *ErrorStatus) *ErrorStatus {
	must.NotNil(err)
	if e, ok := err.(*statusError); ok {
		return e.errorStatus
	}
	if e, ok := err.(grpcStatus); ok {
		return FromStatus(e.GRPCStatus())
	}
	return newError(err)
}

func Canceled(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Canceled, msgAndArgs)
}
func Unknown(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Unknown, msgAndArgs)
}
func InvalidArgument(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.InvalidArgument, msgAndArgs)
}
func DeadlineExceeded(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.DeadlineExceeded, msgAndArgs)
}
func NotFound(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.NotFound, msgAndArgs)
}
func AlreadyExists(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.AlreadyExists, msgAndArgs)
}
func PermissionDenied(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.PermissionDenied, msgAndArgs)
}
func ResourceExhausted(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.ResourceExhausted, msgAndArgs)
}
func FailedPrecondition(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.FailedPrecondition, msgAndArgs)
}
func Aborted(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Aborted, msgAndArgs)
}
func OutOfRange(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.OutOfRange, msgAndArgs)
}
func Unimplemented(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Unimplemented, msgAndArgs)
}
func Internal(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Internal, msgAndArgs)
}
func Unavailable(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Unavailable, msgAndArgs)
}
func DataLoss(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.DataLoss, msgAndArgs)
}
func Unauthenticated(msgAndArgs ...interface{}) *ErrorStatus {
	return newError(codes.Unauthenticated, msgAndArgs)
}

type ErrorStatus struct {
	status status.Status
}

type statusError struct {
	errorStatus *ErrorStatus
}

// impl error interface
func (x *statusError) Error() string {
	return fmt.Sprintf("%s: %s", x.errorStatus.Code(), x.errorStatus.Message())
}

func (x *statusError) GRPCStatus() *status.Status {
	return &x.errorStatus.status
}

// to error interface
func (x *ErrorStatus) Err() error {
	return &statusError{errorStatus: x}
}

func (x *ErrorStatus) Code() codes.Code {
	return x.status.Code()
}
func (x *ErrorStatus) Message() string {
	return x.status.Message()
}
func (x *ErrorStatus) GRPCStatus() *status.Status {
	return &x.status
}

func (x *ErrorStatus) RetryInfo() *errdetails.RetryInfo {
	for _, d := range x.GRPCStatus().Details() {
		if retry, ok := d.(*errdetails.RetryInfo); ok {
			return retry
		}
	}
	return nil
}

func (x *ErrorStatus) WithDebugMessage(debugMsgAndArgs ...interface{}) *ErrorStatus {
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: format("[debug]", debugMsgAndArgs),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *ErrorStatus) WithDetails(details ...proto.Message) *ErrorStatus {
	newStatus, err := x.status.WithDetails(details...)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *ErrorStatus) WithRetryInfo(retryInfo *errdetails.RetryInfo) *ErrorStatus {
	newStatus, err := x.status.WithDetails(retryInfo)
	must.NoError(err)
	return FromStatus(newStatus)
}

func (x *ErrorStatus) WithDebugError(err error) *ErrorStatus {
	must.NotNil(err)
	newStatus, err := x.status.WithDetails(&errdetails.DebugInfo{
		Detail: err.Error(),
	})
	must.NoError(err)
	return FromStatus(newStatus)
}

func newError(code codes.Code, msgAndArgs []interface{}) *ErrorStatus {
	return &ErrorStatus{
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

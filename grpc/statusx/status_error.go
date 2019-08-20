package statusx

import (
	"fmt"
	"google.golang.org/grpc/status"
)

type statusError struct {
	errorStatus *Status
}

// impl error interface
func (x *statusError) Error() string {
	return fmt.Sprintf("%s: %s", x.errorStatus.Code(), x.errorStatus.Message())
}

func (x *statusError) GRPCStatus() *status.Status {
	return &x.errorStatus.status
}

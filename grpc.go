package errs

import (
	"google.golang.org/grpc/status"
)

// GRPCStatus returns a *status.Status representation of *errs.Error
func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.knownCode().GRPC(), e.Error())
}

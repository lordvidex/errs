package status

import "github.com/lordvidex/errs/v2"

// Err converts underlying *errs.Error to *status.Status recommended for grpc handlers
func Err(err error) error {
	instance := errs.Convert(err).(*errs.Error)
	return instance.GRPCStatus().Err()
}

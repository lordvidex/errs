# errs
This package provides a simple error type that can be used to map errors to HTTP and GRPC status codes.

It is heavily inspired by [encore.dev/beta/errs](https://encore.dev/beta/errs) but adds support for GRPC codes.

## Usage
- go docs https://pkg.go.dev/github.com/lordvidex/errs
- check the tests for usage and examples

## Codes
| Code | HTTP Status | GRPC Code | Name |
|------|-------------|-----------|------|
| 0 | http.StatusInternalServerError | codes.Unknown | Unknown Code  |
| 1 | 499 | codes.Canceled | Canceled |
| 2 | http.StatusBadRequest | codes.InvalidArgument | InvalidArgument |
| 3 | http.StatusGatewayTimeout | codes.DeadlineExceeded | DeadlineExceeded |
| 4 | http.StatusUnauthorized | codes.Unauthenticated | Unauthenticated |
| 5 | http.StatusNotFound | codes.NotFound | NotFound |
| 6 | http.StatusConflict | codes.AlreadyExists | AlreadyExists |
| 7 | http.StatusForbidden | codes.PermissionDenied | Forbidden |
| 8 | http.StatusTooManyRequests | codes.ResourceExhausted | ResourceExhausted |
| 9 | http.StatusPreconditionFailed | codes.FailedPrecondition | FailedPrecondition |
| 10 | http.StatusConflict | codes.Aborted | Aborted |
| 11 | httpStatusRequestedRangeNotSatisfiable | codes.OutOfRange | OutOfRange |
| 12 | http.StatusInternalServerError | codes.Internal | Internal |
| 13 | http.StatusServiceUnavailable | codes.Unavailable | Unavailable |
| 14 | http.StatusInternalServerError | codes.DataLoss | DataLoss |

	
## References
- google.golang.org/grpc/codes for grpc codes
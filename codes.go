package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Code is the type that represents an error code.
// It can map to HTTP and gRPC codes.
type Code int

const (
	// Unknown is the default error code.
	Unknown Code = iota

	// Canceled indicates the operation was canceled or unavailable because it was cancelled.
	Canceled

	// InvalidArgument is used when the client sends invalid arguments.
	InvalidArgument

	// DeadlineExceeded means operation expired before completion. This doesn't necessarily mean that the operation failed.
	// It is possible that the operation succeeded but the deadline was exceeded.
	DeadlineExceeded

	// Unauthenticated is used when the client is not authenticated.
	Unauthenticated

	// NotFound is used when the requested resource is not found.
	NotFound

	// AlreadyExists is used when the resource already exists.
	AlreadyExists

	// Forbidden is used when the client is not authorized to perform the requested operation.
	Forbidden

	// ResourceExhausted is used when the client has exhausted some resource.
	ResourceExhausted

	// FailedPrecondition is used when the client sends a request that is not allowed in the current state.
	FailedPrecondition

	// Aborted is used when the client sends a request that cannot be completed due to a conflict e.g. a concurrency issue.
	Aborted

	// OutOfRange means that the operation was attempted past the valid range.
	OutOfRange

	// Internal is used when an internal error occurs.
	Internal

	// Unavailable is used when the service cannot be reached due to some network issues.
	Unavailable

	// DataLoss is used when the service has lost some data.
	DataLoss
)

// CodeSize is the number of codes defined in the STL library.
// All codes defined by default are mapped from 0 to CodeSize - 1.
const CodeSize = 15

// String returns the string representation of the code.
func (c Code) String() string {
	return codeNames[c]
}

// MarshalJSON returns the JSON representation of the code.
func (c Code) MarshalJSON() ([]byte, error) {
	s := c.String()
	return []byte("\"" + s + "\""), nil
}

// HTTP returns the HTTP status code that corresponds to the code.
func (c Code) HTTP() int {
	return httpCodes[c]
}

// GRPC returns the gRPC status code that corresponds to the code.
func (c Code) GRPC() codes.Code {
	return grpcCodes[c]
}

// httpCodes is the mapping of codes to HTTP status codes.
var httpCodes = [...]int{
	Unknown:            http.StatusInternalServerError,
	DataLoss:           http.StatusInternalServerError,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	FailedPrecondition: http.StatusPreconditionFailed,
	ResourceExhausted:  http.StatusTooManyRequests,
	Canceled:           499,
	InvalidArgument:    http.StatusBadRequest,
	Unauthenticated:    http.StatusUnauthorized,
	Forbidden:          http.StatusForbidden,
	NotFound:           http.StatusNotFound,
	AlreadyExists:      http.StatusConflict,
	Internal:           http.StatusInternalServerError,
	OutOfRange:         http.StatusRequestedRangeNotSatisfiable,
	Aborted:            http.StatusConflict,
	Unavailable:        http.StatusServiceUnavailable,
}

// codeNames is the mapping of codes to their string representation.
var codeNames = [...]string{
	Unknown:            "unknown",
	DataLoss:           "data_loss",
	DeadlineExceeded:   "deadline_exceeded",
	FailedPrecondition: "failed_precondition",
	ResourceExhausted:  "resource_exhausted",
	Canceled:           "canceled",
	InvalidArgument:    "invalid_argument",
	Unauthenticated:    "unauthenticated",
	Forbidden:          "forbidden",
	NotFound:           "not_found",
	AlreadyExists:      "already_exists",
	Internal:           "internal",
	OutOfRange:         "out_of_range",
	Aborted:            "aborted",
	Unavailable:        "unavailable",
}

// grpcCodes is the mapping of codes to gRPC codes.
var grpcCodes = [...]codes.Code{
	DataLoss:           codes.DataLoss,
	Unknown:            codes.Unknown,
	InvalidArgument:    codes.InvalidArgument,
	DeadlineExceeded:   codes.DeadlineExceeded,
	FailedPrecondition: codes.FailedPrecondition,
	ResourceExhausted:  codes.ResourceExhausted,
	Canceled:           codes.Canceled,
	Unauthenticated:    codes.Unauthenticated,
	Forbidden:          codes.PermissionDenied,
	NotFound:           codes.NotFound,
	AlreadyExists:      codes.AlreadyExists,
	Internal:           codes.Internal,
	Aborted:            codes.Aborted,
	OutOfRange:         codes.OutOfRange,
	Unavailable:        codes.Unavailable,
}

// Package status is a partial drop-in replacement of grpc/status that works with lordvidex/errs.
// As a drop-in replacement, it simply calls *status.Status functions.
package status

import (
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(c codes.Code, msg string) *status.Status               { return status.New(c, msg) }
func Newf(c codes.Code, format string, a ...any) *status.Status { return status.Newf(c, format, a...) }
func Error(c codes.Code, msg string) error                      { return status.Error(c, msg) }
func Errorf(c codes.Code, format string, a ...any) error        { return status.Errorf(c, format, a...) }
func ErrorProto(s *spb.Status) error                            { return status.ErrorProto(s) }
func FromProto(s *spb.Status) *status.Status                    { return status.FromProto(s) }
func FromError(err error) (s *status.Status, ok bool)           { return status.FromError(err) }
func Convert(err error) *status.Status                          { return status.Convert(err) }
func Code(err error) codes.Code                                 { return status.Code(err) }
func FromContextError(err error) *status.Status                 { return status.FromContextError(err) }

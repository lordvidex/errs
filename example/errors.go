package example

import (
	"github.com/lordvidex/errs"
	"google.golang.org/grpc/codes"
)

const (
	CustomErr errs.Code = iota + errs.CodeSize
	SecondErr
)

func init() {
	Register()
}

func Register() {
	errs.RegisterCode(CustomErr, 400, codes.AlreadyExists, "This is custom error")
	errs.RegisterCode(SecondErr, 399, codes.DataLoss, "This is another error")
}

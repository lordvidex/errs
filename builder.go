package errs

import "fmt"

// B returns a new error builder.
func B() *Builder {
	return &Builder{}
}

type Builder struct {
	code    Code
	msg     string
	details []any
}

// Code sets the code of the error.
func (b *Builder) Code(code Code) *Builder {
	b.code = code
	return b
}

// Msg sets the message of the error.
func (b *Builder) Msg(msg string) *Builder {
	b.msg = msg
	return b
}

// Msgf formats the message using the given format and parameters.
func (b *Builder) Msgf(format string, parameters ...any) *Builder {
	b.msg = fmt.Sprintf(format, parameters...)
	return b
}

// Details adds details to the error.
func (b *Builder) Details(details ...any) *Builder {
	b.details = details
	return b
}

func (b *Builder) Err() error {
	return &Error{
		Code:    b.code,
		Msg:     b.msg,
		Details: b.details,
	}
}

package errs

import (
	"errors"
	"fmt"
	"strings"
)

// B returns a new error builder.
// Optionally, you can pass an error instance to the builder, ONLY the first error will be used.
//
// # Possible cases:
//
// 1. B() -> new error
//
// 2. B(err) -> update existing error (because of pointers) and returns it after calling Err()
//
// 3. B(err1, err2, err3) -> updates only err1 and returns it after calling Err()
func B(initial ...error) *Builder {
	var err *Error
	if len(initial) <= 0 || initial[0] == nil {
		err = new(Error)
	} else {
		err = convert(initial[0])
	}
	return &Builder{err}
}

// WrapB wraps an underlying error and returns a builder for the new error.
func WrapB(err error) *Builder {
	var cause *Error
	errors.As(err, &cause)
	b := B(nil)
	b.err.wrap(cause)
	return b
}

// Builder is used to build an instance of `Error` object.
type Builder struct {
	err *Error
}

// Code sets the code of the error.
func (b *Builder) Code(code Code) *Builder {
	b.err.Code = code
	return b
}

// Msg sets the message of the error.
func (b *Builder) Msg(msg ...string) *Builder {
	b.err.Msg = append(b.err.Msg, cleanStrings(msg)...)
	return b
}

// Op sets the operation where the error occured
func (b *Builder) Op(op string) *Builder {
	b.err.Op = op
	return b
}

func cleanStrings(s []string) []string {
	var r []string
	for _, v := range s {
		if v != "" {
			r = append(r, strings.TrimSpace(v))
		}
	}
	return r
}

// Msgf formats the message using the given format and parameters, similar to fmt.Sprintf.
func (b *Builder) Msgf(format string, parameters ...any) *Builder {
	b.err.Msg = append(b.err.Msg, fmt.Sprintf(format, parameters...))
	return b
}

// Details adds details to the error.
func (b *Builder) Details(details ...any) *Builder {
	b.err.Details = details
	return b
}

// Show sets the show flag of the error.
// If this error is wrapped by another error, it's shown to users.
func (b *Builder) Show() *Builder {
	old := b.err.show
	if old {
		return b
	}

	b.err.shownDepth++
	b.err.show = true
	return b
}

// Err returns new instance of `Error`.
func (b *Builder) Err() error {
	return b.err
}

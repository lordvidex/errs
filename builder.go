package errs

import "fmt"

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
		err = Convert(initial[0]).(*Error)
	}
	return &Builder{err}
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
	b.err.Msg = append(b.err.Msg, rmNilStr(msg)...)
	return b
}

// Op sets the operation where the error occured
func (b *Builder) Op(op string) *Builder {
	b.err.Op = op
	return b
}

func rmNilStr(s []string) []string {
	var r []string
	for _, v := range s {
		if v != "" {
			r = append(r, v)
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

// Err returns new instance of `Error`.
func (b *Builder) Err() error {
	return b.err
}

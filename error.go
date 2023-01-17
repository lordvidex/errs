package errs

import (
	"fmt"
	"strings"
)

type Error struct {
	// Code is the error code of the error. When marshaled to JSON, it will be a string.
	Code Code `json:"code"`

	// Msg is the user-friendly message returned to the client.
	Msg []string `json:"message"`

	// Details is the internal error message returned to the developer.
	Details []any `json:"-"`

	// underlying error
	cause *Error

	// depth of the error tree
	depth int
}

func (e *Error) Error() string {
	return e.Code.String() + ": " + strings.Join(e.Msg, ": ")
}

// Stack returns a description of the error and all it's underlying errors.
func (e *Error) Stack() []string {
	stack := make([]string, e.depth+1)
	for i, err := 0, e; err != nil; err, i = err.cause, i+1 {
		tabOffset := strings.Repeat("\t", i)
		var buf strings.Builder
		write := func(s string) {
			buf.WriteString(tabOffset)
			buf.WriteString(s)
		}
		write(err.Error() + "\n")
		for dx, d := range err.Details {
			write(fmt.Sprintf("\t%d: %v\n", dx, d))
		}
		stack[i] = buf.String()
	}
	return stack
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	if e.cause == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return equalNodes(e, t)
	}
	return false
}

// equalNodes was created because we can't even trust go to compare equality of the error structs.
// Comparison does not involve the underlying errors because we don't want to compare the entire error tree.
//
// The fields considered for equality are error codes and messages. It makes sense to leave details out because two errors
// might be the same but with different details.
func equalNodes(a, b *Error) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Code != b.Code {
		return false
	}
	if len(a.Msg) != len(b.Msg) {
		return false
	}
	for i := range a.Msg {
		if a.Msg[i] != b.Msg[i] {
			return false
		}
	}
	return true
}

// WrapMsg wraps an underlying error with a new error, adding message to the error's previously existing message
func WrapMsg(err error, message ...string) error {
	code := Unknown
	er := Convert(err)
	if er != nil {
		code = er.(*Error).Code
	}
	return WrapCode(er, code, message...)
}

// Wrap wraps an underlying error `child` with a new error `parent`.
//
// - when the child error is nil, the parent error is returned as is.
//
// - when the parent error is nil, the child error is returned as is.
//
// - when both errors are nil, nil is returned.
func Wrap(child, parent error) error {
	parent = Convert(parent)
	child = Convert(child)
	switch {
	case parent == nil && child == nil:
		return nil
	case parent == nil:
		return child
	case child == nil:
		return parent
	default:
		p := parent.(*Error)
		c := child.(*Error)
		p.cause = c
		p.depth = c.depth + 1
		return p
	}
}

// Convert converts any error to an *Error type. If the error is already an *Error, it is returned as is.
// nil errors are returned as nil.
func Convert(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{
		Code: Unknown,
		Msg:  []string{err.Error()},
	}
}

// WrapCode wraps an underlying error with a new error, adding message to the error's previously existing message and setting the error code to code.
func WrapCode(err error, code Code, messages ...string) error {
	er := Convert(err)
	if er == nil {
		return B().Code(code).Msg(rmNilStr(messages)...).Err()
	}
	e := er.(*Error)
	return &Error{
		Code:  code,
		cause: e,
		Msg:   rmNilStr(messages),
		depth: e.depth + 1,
	}
}

package errs

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"
)

// Separator is the default separator between elements of a single error
var Separator = ": "

type Error struct {
	// underlying error
	cause *Error

	// Op operation where error occured
	Op string `json:"op"`

	// Msg is the user-friendly message returned to the client.
	Msg []string `json:"message"`

	// Details is the internal error message returned to the developer.
	Details []any `json:"-"`

	// Code is the error code of the error. When marshaled to JSON, it will be a string.
	Code Code `json:"code"`

	// show is a flag that indicates whether the error would be visible when wrapped by another error
	show bool

	// depth of the error tree
	depth      int
	shownDepth int
}

// Error returns the error in the format "code: message\ninner_code: inner_message" for this error and SHOWN underlying errors.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	buf := strings.Builder{}
	e.writeTo(&buf)

	if e.cause != nil {
		// print the underlying error only if shown
		for inner := range shown(e.cause) {
			buf.WriteString("\n")
			inner.writeTo(&buf)
		}
	}

	return buf.String()
}

// String returns the error in the format "code: message".
func (e *Error) String() string {
	var buf strings.Builder
	e.writeTo(&buf)
	return buf.String()
}

func (e *Error) writeTo(w io.StringWriter) {
	w.WriteString(e.Code.String())
	if len(e.Op) > 0 {
		w.WriteString(Separator + e.Op)
	}

	msgs := strings.Join(cleanStrings(e.Msg), Separator)
	if len(msgs) > 0 {
		w.WriteString(Separator + msgs)
	}
}

// Stack returns a description of the error and all it's underlying errors.
func (e *Error) Stack() string {
	var buf strings.Builder
	for i, er := range all(e) {
		tabOffset := strings.Repeat("\t", i)
		write := func(s string) {
			buf.WriteString(tabOffset)
			buf.WriteString(s)
		}
		write(er.String() + "\n")
		for dx, d := range er.Details {
			write(fmt.Sprintf("\t%d: %v\n", dx, d))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	if e.cause == nil {
		return nil
	}
	return e.cause
}

func (e *Error) wrap(inner *Error) {
	e.cause = inner
	if inner == nil {
		return
	}
	e.depth = inner.depth + 1
	e.shownDepth = inner.shownDepth
	if e.show {
		e.shownDepth++
	}
}

func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
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
	if a.Op != b.Op {
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
		p.wrap(c)
		return p
	}
}

// Convert converts any error to an *Error type. If the error is already an *Error, it is returned as is.
// nil errors are returned as nil.
func Convert(err error) error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{
		Code: Unknown,
		Msg:  []string{err.Error()},
	}
}

// WrapCode wraps an underlying error with a new error, adding message to the error's previously existing message and setting the error code to code.
func WrapCode(err error, code Code, messages ...string) error {
	err = Convert(err)
	if err == nil {
		return B().Code(code).Msg(cleanStrings(messages)...).Err()
	}
	er := err.(*Error)
	e := &Error{
		Code: code,
		Msg:  cleanStrings(messages),
	}
	e.wrap(er)
	return e
}

// shown iterates and yields only errors that should be shown
func shown(e *Error) iter.Seq[*Error] {
	return func(yield func(er *Error) bool) {
		for er := e; er != nil && er.shownDepth > 0; er = er.cause {
			if !er.show {
				continue
			}
			if !yield(er) {
				return
			}
		}
	}
}

// all iterates all inner errors
func all(e *Error) iter.Seq2[int, *Error] {
	return func(yield func(i int, er *Error) bool) {
		for i, er := 0, e; er != nil; i, er = i+1, er.cause {
			if !yield(i, er) {
				return
			}
		}
	}
}

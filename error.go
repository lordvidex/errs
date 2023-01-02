package errs

type Error struct {
	// Code is the error code of the error. When marshaled to JSON, it will be a string.
	Code Code `json:"code"`

	// Msg is the user-friendly message returned to the client.
	Msg string `json:"message"`

	// Details is the internal error message returned to the developer.
	Details []any `json:"-"`
}

func (e *Error) Error() string {
	return e.Code.String() + ": " + e.Msg
}

func Wrap(err error, message string, stacktrace ...any) error {
	return WrapCode(err, message, Unknown, stacktrace...)
}

func WrapCode(err error, message string, code Code, stacktrace ...any) error {
	e := &Error{
		Code:    code,
		Msg:     message,
		Details: stacktrace,
	}
	if err != nil {
		if e2, ok := err.(*Error); ok {
			e.Details = append(e.Details, e2.Details...)
			if e.Code == Unknown {
				e.Code = e2.Code
			}
		} else {
			e.Details = append(e.Details, err.Error())
		}
	}
	return e
}

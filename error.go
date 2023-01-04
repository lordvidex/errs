package errs

type Error struct {
	// Code is the error code of the error. When marshaled to JSON, it will be a string.
	Code Code `json:"code"`

	// Msg is the user-friendly message returned to the client.
	Msg string `json:"message"`

	// Details is the internal error message returned to the developer.
	Details []any `json:"-"`
}

// Error returns the error in the format "code: message".
func (e *Error) Error() string {
	return e.Code.String() + ": " + e.Msg
}

// Wrap returns a new error instance with the given message and stacktrace.
func Wrap(err error, message string, stacktrace ...any) error {
	return WrapCode(err, message, Unknown, stacktrace...)
}

// WrapCode creates a new error instance with the given message, code and stacktrace.
// If the given error is an instance of Error, the code and message will overwrite the ones passed adn the error.Details will be appended to the stacktrace.
func WrapCode(err error, message string, code Code, stacktrace ...any) error {
	e := &Error{
		Code:    code,
		Msg:     message,
		Details: stacktrace,
	}
	if err != nil {
		if e2, ok := err.(*Error); ok {
			e.Details = append(e.Details, e2.Details...)
			// If the code is Unknown, use the code of the given error.
			if e.Code == Unknown {
				e.Code = e2.Code
			}
			// If the message is empty, use the message of the given error.
			if e.Msg == "" {
				e.Msg = e2.Msg
			}
		} else {
			e.Details = append(e.Details, err.Error())
		}
	}
	return e
}

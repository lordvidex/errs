package errs

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrap(t *testing.T) {
	type args struct {
		err        error
		message    string
		code       Code
		stacktrace []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "nil error should still be wrapped with a new error",
			args: args{
				err:        nil,
				message:    "",
				code:       Unknown,
				stacktrace: nil,
			},
			wantErr: &Error{
				Code:    Unknown,
				Msg:     "",
				Details: nil,
			},
		},
		{
			name: "newer error stacktraces should be stacked on top of old ones",
			args: args{
				err:        &Error{Code: NotFound, Msg: "not found", Details: []any{"error1", "error2"}},
				message:    "new error",
				code:       Internal,
				stacktrace: []any{"error3", "error4"},
			},
			wantErr: &Error{
				Code: Internal,
				Msg:  "new error",
				Details: []any{
					"error3",
					"error4",
					"error1",
					"error2",
				},
			},
		},
		{
			name: "wrapped error message should be used when message is empty",
			args: args{
				err: &Error{
					Code: NotFound,
					Msg:  "this message is not found",
				},
				message: "",
				code:    Internal,
			},
			wantErr: &Error{
				Code: Internal,
				Msg:  "this message is not found",
			},
		},
		{
			name: "wrapped error code should be used when code is unknown",
			args: args{
				err: &Error{
					Code:    NotFound,
					Msg:     "not found",
					Details: []any{"error1", "error2"},
				},
				message:    "super message",
				code:       Unknown,
				stacktrace: []any{"error3"},
			},
			wantErr: &Error{
				Code: NotFound,
				Msg:  "super message",
				Details: []any{
					"error3",
					"error1",
					"error2",
				},
			},
		},
		{
			name: "error that does not conform to our Error type should be wrapped with a new error",
			args: args{
				err:        errors.New("not our error"),
				code:       Unavailable,
				message:    "new error",
				stacktrace: []any{"error1", "error2"},
			},
			wantErr: &Error{
				Code:    Unavailable,
				Msg:     "new error",
				Details: []any{"error1", "error2", "not our error"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, WrapCode(tt.args.err, tt.args.message, tt.args.code, tt.args.stacktrace...))
		})
	}
}

func TestError_Wrap(t *testing.T) {
	err := Wrap(nil, "test")
	assert.Equal(t, &Error{Code: Unknown, Msg: "test", Details: nil}, err)
}

func ExampleError_Error() {
	err := B().Code(NotFound).Msg("item not found").Err()
	fmt.Println(err.Error())
	// Output: not_found: item not found
}

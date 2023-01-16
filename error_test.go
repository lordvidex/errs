package errs

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestWrap(t *testing.T) {
	testcases := []struct {
		name          string
		child, parent error
		expect        error
	}{
		{
			name:   "nil child",
			child:  nil,
			parent: errors.New("parent"),
			expect: &Error{
				Msg: []string{"parent"},
			},
		},
		{
			name:   "nil parent",
			child:  errors.New("child"),
			parent: nil,
			expect: &Error{
				Msg: []string{"child"},
			},
		},
		{
			name:   "both nil",
			child:  nil,
			parent: nil,
			expect: nil,
		},
		{
			name:   "child and parent are provided",
			child:  B().Code(NotFound).Msg("child").Err(),
			parent: B().Code(Internal).Msg("parent").Err(),
			expect: &Error{
				Code: Internal,
				Msg:  []string{"parent"},
				cause: &Error{
					Code: NotFound,
					Msg:  []string{"child"},
				},
				depth: 1,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, Wrap(tc.child, tc.parent))
		})
	}
}

func TestWrapMsg(t *testing.T) {
	type args struct {
		err     error
		message string
		code    Code
	}
	tests := []struct {
		name    string
		args    args
		wantErr func(args) error
	}{
		{
			name: "nil error should still be wrapped with a new error",
			args: args{
				err:     nil,
				message: "test",
				code:    Unknown,
			},
			wantErr: func(_ args) error {
				return &Error{
					Code:    Unknown,
					Msg:     []string{"test"},
					Details: nil,
				}
			},
		},
		{
			name: "error should be wrapped with a new error",
			args: args{
				err:     &Error{Code: NotFound, Msg: []string{"not found"}, Details: []any{"error1", "error2"}},
				message: "new error",
				code:    Internal,
			},
			wantErr: func(x args) error {
				return &Error{
					Code:  Internal,
					Msg:   []string{"new error"},
					cause: x.err.(*Error),
					depth: 1,
				}
			},
		},
		{
			name: "empty message should be ignored",
			args: args{
				err: &Error{
					Code: NotFound,
					Msg:  []string{"this message is not found"},
				},
				message: "",
				code:    Internal,
			},
			wantErr: func(e args) error {
				return &Error{
					cause: e.err.(*Error),
					depth: 1,
					Code:  Internal,
				}
			},
		},
		{
			name: "error that does not conform to our Error type should be wrapped with a new error",
			args: args{
				err:     errors.New("not our error"),
				code:    Unavailable,
				message: "new error",
			},
			wantErr: func(_ args) error {
				return &Error{
					Code:  Unavailable,
					Msg:   []string{"new error"},
					depth: 1,
					cause: B().Msg("not our error").Err().(*Error),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr(tt.args), WrapCode(tt.args.err, tt.args.code, tt.args.message))
		})
	}
}

func ExampleError_Stack() {
	err1 := B().Code(NotFound).Msg("item not found").Err()
	err2 := B().Code(Internal).Msg("internal error").Err()
	err, _ := Wrap(err2, err1).(*Error)
	fmt.Print(err.Stack())
}

func TestError_Wrap(t *testing.T) {
	err := WrapMsg(nil, "test")
	assert.Equal(t, &Error{Code: Unknown, Msg: []string{"test"}, Details: nil}, err)
}

func ExampleError_Error() {
	err := B().Code(NotFound).Msg("item not found").Err()
	fmt.Println(err.Error())
	// Output: not_found: item not found
}

var e error

func BenchmarkWrapCode(b *testing.B) {
	var _e error
	for i := 0; i < b.N; i++ {
		_e = WrapCode(nil, NotFound, "item not found", strconv.Itoa(i))
	}
	e = _e
}

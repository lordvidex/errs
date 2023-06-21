package errs

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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
			parent: B().Code(Internal).Msg("parent").Op("UserRepository.CreateUser").Err(),
			expect: &Error{
				Code: Internal,
				Msg:  []string{"parent"},
				Op:   "UserRepository.CreateUser",
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
	// this test contains only events where the error is non-nil, as other cases are handled by WrapCode
	t.Run("non nil error", func(t *testing.T) {
		initialError := B().Code(NotFound).Msg("not found").Err()
		casted := initialError.(*Error)
		err := WrapMsg(initialError, "new error")
		assert.Equal(t, &Error{
			Code:  NotFound,
			Msg:   []string{"new error"},
			cause: casted,
			depth: 1,
		}, err)
	})
}

// TestErrorIs is useful to test if the Unwrap function is implemented correctly and also to make sure that clients
// can use errors.Is to check for existence of some specific error even if it is underlying.
func TestErrorIs(t *testing.T) {
	makeNested := func(digits ...int) error {
		if len(digits) == 0 {
			return nil
		}
		var err error
		for i := len(digits) - 1; i >= 0; i-- {
			err = Wrap(err, B().Code(InvalidArgument).Msgf("invalid argument %d", digits[i]).Err())
		}
		return err
	}
	testcases := []struct {
		name   string
		err    error
		target error
		expect bool
	}{
		{
			name:   "nil errors",
			err:    nil,
			target: nil,
			expect: true,
		},
		{
			name: "nil target",
			err:  B().Code(NotFound).Msg("not found").Err(),
			// target is nil, so it should always return false
			target: nil,
			expect: false,
		},
		{
			name:   "nil error",
			err:    nil,
			target: B().Code(NotFound).Msg("not found").Err(),
			expect: false,
		},
		{
			name:   "same error",
			err:    B().Code(NotFound).Msg("not found").Err(),
			target: B().Code(NotFound).Msg("not found").Err(),
			expect: true,
		},
		{
			name:   "nested errors",
			err:    makeNested(1, 2, 3, 4),
			target: makeNested(3),
			expect: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, errors.Is(tc.err, tc.target))
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	err := WrapMsg(nil, "test")
	err2 := Wrap(err, B().Code(NotFound).Err())
	assert.Nil(t, errors.Unwrap(err))
	assert.Equal(t, err, errors.Unwrap(err2))
}

func TestError_Is(t *testing.T) {
	myerr := B().Code(NotFound).Msg("item not found").Err()
	cases := []struct {
		name   string
		err    error
		expect bool
	}{
		{"nil error", nil, false},
		{"does not conform to class", errors.New("test"), false},
		{"conforms to class and equal", myerr, true},
		{"not equal code", B().Code(InvalidArgument).Err(), false},
		{"not equal messages - message not equal", B().Code(NotFound).Msg("not found").Err(), false},
		{"not equal messages - first equals, second not", B().Code(NotFound).Msg("item not found", "not found").Err(), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, myerr.(*Error).Is(tc.err))
		})
	}
}

// TestError_equal is a private test.
// equalNodes cannot be used outside the package and thus this test is written to fulfill all righteousness and make sure
// the function is working as expected.
func TestError_equal(t *testing.T) {
	cases := []struct {
		name   string
		a, b   *Error
		expect bool
	}{
		{"both nil", nil, nil, true},
		{"a nil", nil, new(Error), false},
		{"b nil", new(Error), nil, false},
		{"both not nil", new(Error), new(Error), true},
		{"both not nil - different", &Error{Code: NotFound}, &Error{Code: InvalidArgument}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, equalNodes(tc.a, tc.b))
		})
	}
}

func TestWrapCode(t *testing.T) {
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
	err := B().Code(NotFound).Msg("item not found").Op("FetchItem").Err()
	fmt.Println(err.Error())
	// Output: not_found: FetchItem: item not found
}

var e error

func BenchmarkWrapCode(b *testing.B) {
	var _e error
	for i := 0; i < b.N; i++ {
		_e = WrapCode(nil, NotFound, "item not found", strconv.Itoa(i))
	}
	e = _e
}

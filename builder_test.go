package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleBuilder() {
	b := B().Code(Unknown).Msg("unknown error").Details("details")
	fmt.Println(b.Err())

	// Output: unknown: unknown error
}

func ExampleBuilder_Msgf() {
	b := B().Code(NotFound).Msgf("file not found: %s, %d", "details", 123).Op("File.Open")
	err := b.Err()
	str, _ := json.Marshal(err)
	fmt.Println(string(str))

	// Output: {"op":"File.Open","message":["file not found: details, 123"],"code":"not_found"}
}

func TestWrapB(t *testing.T) {
	tests := []struct {
		// given
		inner            error
		name             string
		expectDepth      int
		expectShownDepth int
		expectErr        string
	}{
		{
			name:      "wrapping nil error",
			inner:     nil,
			expectErr: "unknown",
		},
		{
			name:             "wrapping existing wrapped error",
			inner:            Wrap(B().Msg("inner0").Show().Err(), B().Msg("inner1").Err()),
			expectDepth:      2,
			expectShownDepth: 1,
			expectErr:        "unknown\nunknown: inner0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapB(tt.inner)
			assert.Equal(t, tt.expectDepth, got.err.depth)
			assert.Equal(t, tt.expectShownDepth, got.err.shownDepth)
			assert.Equal(t, tt.expectErr, got.err.Error())
		})
	}
}

func TestBuilder(t *testing.T) {
	t.Run("should update existing error", func(t *testing.T) {
		err := B().Code(NotFound).Msg("item not found").Err() // err is updated by pointer
		updated := B(err).Code(InvalidArgument).Msg("invalid argument").Err()
		assert.Equal(t, err, updated)
	})
	t.Run("passed nil error", func(t *testing.T) {
		err := B(nil).Code(NotFound).Msg("item not found").Err()
		assert.Equal(t, &Error{Code: NotFound, Msg: []string{"item not found"}}, err)
	})
	t.Run("passed error that does not conform to our Error type", func(t *testing.T) {
		err := errors.New("not our error")
		updated := B(err).Code(NotFound).Msg("item not found").Err()
		assert.NotEqual(t, err, updated) // since it does not conform to our error, it must have been converted
		assert.Equal(
			t,
			&Error{Code: NotFound,
				Msg: []string{"not our error", "item not found"},
			},
			updated,
		)
	})
}

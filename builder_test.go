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
	err := b.Err()
	str, _ := json.Marshal(err)
	fmt.Println(string(str))

	// Output: {"code":"unknown","message":["unknown error"],"op":""}
}

func ExampleBuilder_Msgf() {
	b := B().Code(NotFound).Msgf("file not found: %s, %d", "details", 123).Op("File.Open")
	err := b.Err()
	str, _ := json.Marshal(err)
	fmt.Println(string(str))

	// Output: {"code":"not_found","message":["file not found: details, 123"],"op":"File.Open"}
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

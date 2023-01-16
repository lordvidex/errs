package errs

import (
	"encoding/json"
	"fmt"
)

func ExampleBuilder() {
	b := B().Code(Unknown).Msg("unknown error").Details("details")
	err := b.Err()
	str, _ := json.Marshal(err)
	fmt.Println(string(str))

	// Output: {"code":"unknown","message":["unknown error"]}
}

func ExampleBuilder_Msgf() {
	b := B().Code(NotFound).Msgf("file not found: %s, %d", "details", 123)
	err := b.Err()
	str, _ := json.Marshal(err)
	fmt.Println(string(str))

	// Output: {"code":"not_found","message":["file not found: details, 123"]}
}

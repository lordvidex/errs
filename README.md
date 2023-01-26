# errs [![CI Test](https://github.com/lordvidex/errs/actions/workflows/test.yml/badge.svg)](https://github.com/lordvidex/errs/actions/workflows/test.yml)
This package provides a simple error type that can be used to map errors to HTTP and GRPC status codes.

It is heavily inspired by [encore.dev/beta/errs](https://encore.dev/beta/errs) but adds support for GRPC codes.

## Usage
- go docs https://pkg.go.dev/github.com/lordvidex/errs
- Examples:

```go
package main

import (
	"fmt"
	"github.com/lordvidex/errs"
	"encoding/json"
)

func main() {
	err := errs.B().Code(errs.NotFound).Msg("user not found").Err()
	// err.Error() == "user not found"
	// err.HTTPCode() == 404
	// err.GRPCCode() == 5
	b, _ := json.Marshal(err)
	fmt.Println(string(b)) 
	// Outputs: {"code":"not_found", "message":"user not found"}
}
```
- check the tests for more usage and examples

## References
- google.golang.org/grpc/codes for grpc codes

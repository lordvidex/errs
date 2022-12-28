package errs

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestCode_HTTP(t *testing.T) {
	for i := 0; i < CodeSize; i++ {
		testName := fmt.Sprintf("%s should not be zero", Code(i).String())
		t.Run(testName, func(t *testing.T) {
			if Code(i).HTTP() == 0 {
				t.Errorf("Code(%d).HTTP() == 0", i)
			}
		})
	}
}

func TestCode_GRPC(t *testing.T) {
	for i := 0; i < CodeSize; i++ {
		testName := fmt.Sprintf("%s should not be zero", Code(i).String())
		t.Run(testName, func(t *testing.T) {
			if Code(i).GRPC() == codes.OK {
				t.Errorf("Code(%d).GRPC() == 0", i)
			}
		})
	}
}

func TestCode_String(t *testing.T) {
	for i := 0; i < CodeSize; i++ {
		testName := fmt.Sprintf("%s should not be empty", Code(i).String())
		t.Run(testName, func(t *testing.T) {
			if Code(i).String() == "" {
				t.Errorf("Code(%d).String() == \"\"", i)
			}
		})
	}
}

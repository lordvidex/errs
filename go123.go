//go:build go1.23
// +build go1.23

package errs

import "iter"

// shown iterates and yields only errors that should be shown
// iter.Seq2 is used for compatibility with slice ranging.
// After update to >=1.23, change to iter.Seq
func shown(e *Error) iter.Seq2[int, *Error] {
	return func(yield func(int, *Error) bool) {
		for i, er := 0, e; er != nil && er.shownDepth > 0; i, er = i+1, er.cause {
			if !er.show {
				continue
			}
			if !yield(i, er) {
				return
			}
		}
	}
}

// all iterates all inner errors
func all(e *Error) iter.Seq2[int, *Error] {
	return func(yield func(i int, er *Error) bool) {
		for i, er := 0, e; er != nil; i, er = i+1, er.cause {
			if !yield(i, er) {
				return
			}
		}
	}
}

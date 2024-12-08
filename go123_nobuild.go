//go:build !go1.23
// +build !go1.23

package errs

func shown(e *Error) []*Error {
	var arr []*Error
	for er := e; er != nil; er = er.cause {
		if !er.show {
			continue
		}
		arr = append(arr, er)
	}
	return arr
}

func all(e *Error) []*Error {
	arr := []*Error{e}
	for i, er := 0, e.cause; er != nil; i, er = i+1, er.cause {
		arr = append(arr, er)
	}
	return arr
}

package testpkg

import "errors"

func fn1() {
	err := errors.New("x") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		_ = 1
	}
}

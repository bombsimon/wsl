package testpkg

import "errors"

func fn1() {
	err := errors.New("x")
	if err != nil { // want "unnecessary whitespace decreases readability"
		panic(err)
	}
}

func fn2() {
	one := 1
	err := errors.New("x")
	if err != nil { // want "unnecessary whitespace decreases readability"
		panic(err)
	}

	_ = 1
}

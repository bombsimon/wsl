package testpkg

import "errors"

func fn1() {
	one := 1
	two := 2 // want "missing whitespace decreases readability"
	if two == 2 {
		panic(err)
	}
}

func fn2() {
	two := 2
	one := 1
	if two == 2 { // want "missing whitespace decreases readability"
		panic(err)
	}
}

func fn3() {
	err := errors.New("error")

	if err != nil {
		panic(err)
	}
}
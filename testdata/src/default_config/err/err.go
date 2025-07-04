package testpkg

import "errors"

func fn1() {
	err := errors.New("x") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		panic(err)
	}

	err2 := errors.New("y") // want +1 `unnecessary whitespace \(err\)`

	if err2 == nil {
		panic(err)
	}
}

func fn11() { // want +2 `unnecessary whitespace \(err\)`
	err := errors.New("x")

	if err != nil {
		panic(err)
	}
}

func fn12() {
	err := errors.New("x")
	// Some comment
	if err != nil {
		panic(err)
	}
}

func fn13() {
	err := errors.New("x")

	// Some comment
	if err != nil {
		panic(err)
	}
}

func fn2() {
	a := 1
	err := errors.New("x") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		panic(err)
	}

	_ = a
}

func fn21() {
	a := 1 // want +2 `unnecessary whitespace \(err\)`
	err := errors.New("x")

	if err != nil {
		panic(err)
	}

	_ = a
}

func fn3() {
	a := 1
	err := errors.New("x") // want `missing whitespace above this line \(too many statements above if\)`
	if err != nil {
		panic(err)
	}

	_ = a
}

func fn4() {
	msg := "not an error"
	err := &msg

	if err != nil {
		panic(err)
	}
}

func fn5() {
	err := errors.New("x") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		panic(err)
	}

	if false {
		_ = 1
	}
}

func fn6() {
	err := errors.New("x")

	if false {
		panic(err)
	}
}

func fn7() {
	err, ok := errors.New("x"), true

	if !ok {
		panic(err)
	}
}

func fn8() {
	aErr := errors.New("a")
	bErr := errors.New("b")

	if aErr != nil && bErr != nil {
		panic(aErr)
	}
}

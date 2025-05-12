package testpkg

import "fmt"

func Call(fn func()) func() {
	return fn
}

func fn0() { fmt.Println("Hello, World") }

func fn1() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	fmt.Println("Hello, World")
}

func fn2() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	// Space before comment.
	fmt.Println("Hello, World")
}

func fn3() { // want +2 `unnecessary whitespace \(leading-whitespace\)`
	// Space after comment

	fmt.Println("Hello, World")
}

func fn4() {
	fmt.Println("Hello, World")
}

func fn5() {
	// Comment without space before or after.
	fmt.Println("Hello, World")
}

func fn51() {
	// Comment with

	// a newline between
	_ = 1
}

func fn6() {
	if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}
}

func fn7() {
	if true {
		_ = 1
	} else if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	} else { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}
}

func fn8(a string, b any, s []string) { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	} else if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}

	for i := 0; i < 1; i++ { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}

	for n := range []int{} { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = n
	}

	for range s { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}

	switch a { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	case "a":
	}

	switch b.(type) { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	case int:
	}

	f := func() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}

	f2 := Call(func() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	})

	_ = f
	_ = f2
}

func fn9() {
	switch {
	case 1: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1

	case 2: // want +2 `unnecessary whitespace \(leading-whitespace\)`
		// This is a comment

		_ = 2

	default: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 3

	}
}

func fn10() {
	switch {
	case true || false:
		fmt.Println("ok")
	case true, false:
		fmt.Println("ok")
	case true || false:
		fmt.Println("ok")
	case true,
		false: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		fmt.Println("nok")
	case true ||
		false: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		fmt.Println("nok")
	case true,
		false:
		fmt.Println("ok")

	case true ||
		false:
		fmt.Println("ok")

	case true, false:
		fmt.Println("all")
	}
}

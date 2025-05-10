package testpkg

import "fmt"

func Call(fn func()) func() {
	return fn
}

func fn1() {
	fmt.Println("Hello, World") // want +1 `unnecessary whitespace \(trailing-whitespace\)`

}

func fn2() {
	fmt.Println("Hello, World")  // want +2 `unnecessary whitespace \(trailing-whitespace\)`
	// Comment with wihtespace

}

func fn3() {
	fmt.Println("Hello, World")
	// Comment without space before or after.
}

func fn4() {
	_ = 1
	// Comment with

	// a newline between
}

func fn5() {
	if true {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}
}

func fn6() {
	if true {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	} else if true {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	} else {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}
}

func fn8(a string, b any, s []string) {
	if true {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	} else if true {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}

	for i := 0; i < 1; i++ {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}

	for n := range []int{} {
		_ = n // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}

	for range s {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}

	switch a {
	case "a":

	}

	switch b.(type) {
	case int:

	}

	f := func() {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	}

	f2 := Call(func() {
		_ = 1 // want +1 `unnecessary whitespace \(trailing-whitespace\)`

	})

	_ = f
	_ = f2
}

package testpkg

import "fmt"

func Call(fn func()) func() {
	return fn
}

func fn1() {
	fmt.Println("Hello, World")

} // want "unnecessary whitespace decreases readability"

func fn2() {
	fmt.Println("Hello, World")
	// Comment with wihtespace

} // want "unnecessary whitespace decreases readability"

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
		_ = 1

	} // want "unnecessary whitespace decreases readability"
}

func fn6() {
	if true {
		_ = 1

	} else if true { // want "unnecessary whitespace decreases readability"
		_ = 1

	} else { // want "unnecessary whitespace decreases readability"
		_ = 1

	} // want "unnecessary whitespace decreases readability"
}

func fn8(a string, b any, s []string) {
	if true {
		_ = 1

	} else if true { // want "unnecessary whitespace decreases readability"
		_ = 1

	} // want "unnecessary whitespace decreases readability"

	for i := 0; i < 1; i++ {
		_ = 1

	} // want "unnecessary whitespace decreases readability"

	for n := range []int{} {
		_ = n

	} // want "unnecessary whitespace decreases readability"

	for range s {
		_ = 1

	} // want "unnecessary whitespace decreases readability"

	switch a {
	case "a":

	}

	switch b.(type) {
	case int:

	}

	f := func() {
		_ = 1

	} // want "unnecessary whitespace decreases readability"

	f2 := Call(func() {
		_ = 1

	}) // want "unnecessary whitespace decreases readability"

	_ = f
	_ = f2
}

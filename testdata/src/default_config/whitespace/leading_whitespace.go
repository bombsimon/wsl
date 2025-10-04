package testpkg

import "fmt"

func CallLeading(fn func()) func() {
	return fn
}

func oneliner() { fmt.Println("Hello, World") }

func leadingWhitespace() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	fmt.Println("Hello, World")
}

func leadingWhitespaceThenComment() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	// Space before comment.
	fmt.Println("Hello, World")
}

func commentThenLeadingWhitespace() { // want +2 `unnecessary whitespace \(leading-whitespace\)`
	// Space after comment

	fmt.Println("Hello, World")
}

func noWhitespace() {
	fmt.Println("Hello, World")
}

func leadingComment() {
	// Comment without space before or after.
	fmt.Println("Hello, World")
}

func multipleLeadingComments() {
	// Comment with

	// a newline between
	_ = 1
}

func leadingIf() {
	if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}
}

func leadingElseElseIf() {
	if true {
		_ = 1
	} else if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	} else { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	}
}

func leadingFnAndIfChain(a string, b any, s []string) { // want +1 `unnecessary whitespace \(leading-whitespace\)`

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

	f2 := CallLeading(func() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1
	})

	_ = f
	_ = f2
}

func leadingSwitch() {
	switch 1 {
	case 1: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 1

	case 2: // want +2 `unnecessary whitespace \(leading-whitespace\)`
		// This is a comment

		_ = 2

	default: // want +1 `unnecessary whitespace \(leading-whitespace\)`

		_ = 3

	}
}

func leadingComplexCase() {
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

func leadingStructAssignment() {
	aBool := func(f func() bool) bool {
		return f()
	}

	type T struct {
		B  bool
		Fn func()
	}

	t := T{
		B: aBool(func() bool { // want +1 `unnecessary whitespace \(leading-whitespace\)`

			return true
		}),
		Fn: func() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

			fmt.Println("nok")
		},
	}

	_ = t
}

func variousLeading() {
	// Assorted weird anonymous functions inside subexpressions,
	// unlikely to be encountered in practice.

	// parenthesized expression
	_ = (func() bool { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		return true
	})

	// binary op
	_ = 1 + func() int { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		return 123
	}()

	// type assertion expression
	_ = func() interface{} { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		return 2
	}().(int)

	// test for proper FuncType traversal in type assertion below
	var f any = func() { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		panic("TODO")
	}

	_ = f.(func())

	var arr [10]byte

	// slice expression
	_ = arr[func() int { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		return 5
	}()]
}

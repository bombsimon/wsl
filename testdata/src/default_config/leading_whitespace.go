package testpkg

import "fmt"

func fn1() { // want "unnecessary whitespace decreases readability"

	fmt.Println("Hello, World")
}

func fn2() { // want "unnecessary whitespace decreases readability"

	// Space before comment.
	fmt.Println("Hello, World")
}

func fn3() {
	// Space after comment // want "unnecessary whitespace decreases readability"

	fmt.Println("Hello, World")
}

func fn4() {
	fmt.Println("Hello, World")
}

func fn5() {
	// Comment without space before or after.
	fmt.Println("Hello, World")
}

func fn6() {
	if true { // want "unnecessary whitespace decreases readability"

		_ = 1
	}
}

func fn7() {
	if true {
		_ = 1
	} else if true { // want "unnecessary whitespace decreases readability"

		_ = 1
	} else { // want "unnecessary whitespace decreases readability"

		_ = 1
	}
}

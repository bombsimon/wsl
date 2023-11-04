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
	// Space after comment - v5.0.0 BREAKING CHANGE

	fmt.Println("Hello, World")
}

func fn4() {
	fmt.Println("Hello, World")
}

func fn5() {
	// Comment without space before or after.
	fmt.Println("Hello, World")
}

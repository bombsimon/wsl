package pkg

import "fmt"

func Fn() {
	// Comment one

	// Comment two
	fmt.Println("a")
}

func Fn2() { // want "block should not start with a whitespace"
	// Comment one

	// Comment two

	fmt.Println("a")
}

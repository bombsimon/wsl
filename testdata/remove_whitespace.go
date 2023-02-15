package testpkg

import "fmt"

func RemoveWhitespaceNoComments() { // want "block should not start with a whitespace"

	a := 1
	if a < 2 { // want "block should not start with a whitespace"

		a = 2

	} // want "block should not end with a whitespace"

	switch { // want "block should not start with a whitespace"

	case true:
		fmt.Println("true")

	default:
		fmt.Println("false")

	} // want "block should not end with a whitespace"

	_ = a

} // want "block should not end with a whitespace"

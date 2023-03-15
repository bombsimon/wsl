package pkg

import "fmt"

func Fn() {
	if true {
		fmt.Println("a")
		// Comment
	}

	if true {
		fmt.Println("a")
		// Comment

		// Comment
	} // want "block should not end with a whitespace"
}

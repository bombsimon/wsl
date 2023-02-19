package pkg

import "fmt"

func Fn() {
	if true {
		fmt.Println("a")
		// Comment
	} // want "block should not end with a whitespace"
}

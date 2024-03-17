package testpkg

import "fmt"

func fn1() {
	one := 1 // want "missing whitespace decreases readability"

	if 1 == 1 {
		fmt.Println(1)
	} else if 1 == 2 {
		fmt.Println(2)
	} else {
		fmt.Println(2)
	}

	_ = one
}

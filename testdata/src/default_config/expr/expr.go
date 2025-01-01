package testpkg

import "fmt"

func fn() {
	a := 1
	b := 2
	c := 3
	a = 4
	b = 3

	_ = a
	_ = b
	_ = c
}

func fn2() {
	a := 1
	b := 2
	fmt.Println("") // want "missing whitespace decreases readability"
	c := 3          // want "missing whitespace decreases readability"
	d := 4

	_ = a
	_ = b
	_ = c
	_ = d
}

func fn3() {
	a := 1
	b := 2
	fmt.Println(b)
	c := 3 // want "missing whitespace decreases readability"
	d := 4

	_ = a
	_ = b
	_ = c
	_ = d
}

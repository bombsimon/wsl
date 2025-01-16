package testpkg

func fn() {
	a := 1
	b := 2
	c := 3

	a = 1
	b = 1
	c = 1

	d := 1
	a = 3  // want "missing whitespace decreases readability"
	e := 4 // want "missing whitespace decreases readability"
}

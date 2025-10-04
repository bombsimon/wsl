package testpkg

func fn() {
	a := 1
	b := 2
	c := 3

	a = 1
	b = 1
	c = 1

	d := 1
	a = 3  // want `missing whitespace above this line \(invalid statement above assign\)`
	e := 4 // want `missing whitespace above this line \(invalid statement above assign\)`

	_, _, _, _, _ = a, b, c, d, e
}

package testpkg

func fn1() int {
	_ = 1
	_ = 2
	_ = 3
	_ = 4
	return 1
}

func fn2() int {
	_ = 1
	_ = 2
	_ = 3
	_ = 4
	_ = 5
	return 1 // want `missing whitespace above this line \(too many lines above return\)`
}

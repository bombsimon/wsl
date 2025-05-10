package testpkg

func fn1() int {
	return 1
}

func fn2() int {
	_ = 1
	return 1
}

func fn3() int {
	_ = 1
	_ = 2
	return 1 // want `missing whitespace above this line \(too many lines above return\)`
}

func fn4() int {
	if true {
		_ = 1
	}
	return 1 // want `missing whitespace above this line \(too many lines above return\)`
}

func fn5() int {
	if true {
		_ = 1
		_ = 2
		return 1 // want `missing whitespace above this line \(too many lines above return\)`
	}
	return 1 // want `missing whitespace above this line \(too many lines above return\)`
}

func fn6() func() {
	_ = 1
	_ = 2
	return func() { // want `missing whitespace above this line \(too many lines above return\)`
		_ = 1
	}
}

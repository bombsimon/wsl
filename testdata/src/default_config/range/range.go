package testpkg

func fn1() {
	a := []int{}
	b := []int{} // want `missing whitespace above this line \(too many statements above range\)`
	for range b {
		panic(1)
	}

	_, _ = a, b
}

func fn2() {
	b := []int{}
	a := []int{}
	for range b { // want `missing whitespace above this line \(no shared variables above range\)`
		panic(1)
	}

	_, _ = a, b
}

func fn3() {
	for range make([]int, 0) {
		panic("")
	}
	for range make([]int, 0) { // want `missing whitespace above this line \(invalid statement above range\)`
		panic("")
	}
}

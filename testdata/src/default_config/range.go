package testpkg

func fn1() {
	one := []int{}
	two := []int{} // want "missing whitespace decreases readability"
	for range two {
		panic(err)
	}
}

func fn2() {
	two := []int{}
	one := []int{}
	for range two { // want "missing whitespace decreases readability"
		panic(err)
	}
}

func fn3() {
	for range make([]int{}, 0) {
		panic(foo)
	}
	for range make([]int{}, 0) { // want "missing whitespace decreases readability"
		panic(bar)
	}
}

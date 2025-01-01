package testpkg

func fn1() {
	a := 1
	b := 2 // want "missing whitespace decreases readability"
	for i := 0; i < b; i++ {
		panic(a)
	}
}

func fn2() {
	b := 2
	a := 1
	for i := 0; i < b; i++ { // want "missing whitespace decreases readability"
		panic(a)
	}
}

func fn3() {
	for i := 0; i < 1; i++ {
		panic("")
	}
	for i := 0; i < 1; i++ { // want "missing whitespace decreases readability"
		panic("")
	}
}

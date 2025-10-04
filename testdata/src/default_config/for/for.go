package testpkg

func keepCuddledVar() {
	a := 1
	b := 2 // want `missing whitespace above this line \(too many statements above for\)`
	for i := 0; i < b; i++ {
		panic(1)
	}

	_, _ = a, b
}

func keepNone() {
	b := 2
	a := 1
	for i := 0; i < b; i++ { // want `missing whitespace above this line \(no shared variables above for\)`
		panic(1)
	}

	_, _ = a, b
}

func separateBlocks() {
	for i := 0; i < 1; i++ {
		panic("")
	}
	for i := 0; i < 1; i++ { // want `missing whitespace above this line \(invalid statement above for\)`
		panic("")
	}
}

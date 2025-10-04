package testpkg

func Fn(_ int) {}

func fn1() {
	x := 1
	if true {
		y, z := 1, 2
		if false {
			z++
			if true { // want `missing whitespace above this line \(no shared variables above if\)`
				y++
				x++
			}
		}
	}
}

func fn2() {
	a := 1
	b := 1
	c := 1

	a = 1
	if true {
		a++
	} else if false {
		b++
	} else {
		c++
	}

	b = 1
	if true {
		a++
	} else if false {
		b++
	} else {
		c++
	}

	c = 1
	if true {
		a++
	} else if false {
		b++
	} else {
		c++
	}
}

func fn3(a, b, c int) {
	a = 1
	for range make([]int, 10) {
		if true {
			c++

			Fn(a)
		}
	}

	b = 1
	for range make([]int, 10) { // want `missing whitespace above this line \(no shared variables above range\)`
		if true {
			c++

			Fn(a)
		}
	}

	a = 1
	for i := 0; i < 3; i++ {
		if true {
			c++

			Fn(a)
		}
	}

	b = 1
	for i := 0; i < 3; i++ { // want `missing whitespace above this line \(no shared variables above for\)`
		if true {
			c++

			Fn(a)
		}
	}
}

func fn4(a, b, c int) {
	a = 1
	switch {
	case true:
		c++
		a++
	case false:
		c++
		a--
	}

	b = 1
	switch { // want `missing whitespace above this line \(no shared variables above switch\)`
	case true:
		c++
		a++
	case false:
		c++
		a--
	}
}

func fn5(a, b, c int) {
	a = 1
	go func() {
		c++
		a++
	}()

	b = 1
	go func() { // want `missing whitespace above this line \(no shared variables above go\)`
		c++
		a++
	}()
}

func fn6(a, b, c int) {
	a = 1
	defer func() {
		if true {
			c++
			a++
		}
	}()

	b = 1
	defer func() { // want `missing whitespace above this line \(no shared variables above defer\)`
		c++
		a++
	}()
}

package testpkg

func fn1() {
	one := 1
	two := 2 // want "missing whitespace decreases readability"
	switch two {
	case 1:
	case 2:
	case 3:
		panic(err)
	}
}

func fn2() {
	two := 2
	one := 1
	switch two { // want "missing whitespace decreases readability"
	case 1:
	case 2:
	case 3:
		panic(err)
	}
}

func fn3() {
	switch true {
	case true:
	case false:
	}
	switch true { // want "missing whitespace decreases readability"
	case true:
	case false:
		panic(err)
	}
}

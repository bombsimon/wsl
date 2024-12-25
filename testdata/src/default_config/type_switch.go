package testpkg

func fn1() {
	var v any

	v = 1
	switch v.(type) {
	case int:
	case string:
	}
}

func fn2() {
	var v any
	v = 1

	one := 1
	switch v.(type) { // want "missing whitespace decreases readability"
	case int:
	case string:
	}
}

func fn3(in any) {
	if true {
		return
	}

	switch v := in.(type) {
	case string:
		return
	default:
	}
}

func fn4(in any) {
	one := getOne()
	two := getTwo() // want "missing whitespace decreases readability"
	switch two.(type) {
	case string:
		return
	default:
	}
}

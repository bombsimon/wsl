package testpkg

func Fn() any {
	return nil
}

func fn1() {
	var v any

	v = 1
	switch v.(type) {
	case int:
	case string:
	}
}

func fn2() {
	var a any = 1

	b := 1
	switch a.(type) { // want `missing whitespace above this line \(no shared variables above type-switch\)`
	case int:
	case string:
	}

	_ = b
}

func fn3(in any) {
	if true {
		return
	}

	switch a := in.(type) {
	case string:
		return
	default:
		_ = a
	}
}

func fn4(in any) {
	a := Fn()
	b := Fn() // want `missing whitespace above this line \(too many statements above type-switch\)`
	switch b.(type) {
	case string:
		return
	default:
		_ = a
	}
}

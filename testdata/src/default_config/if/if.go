package testpkg

import "errors"

func fn1() {
	a := 1
	b := 2 // want `missing whitespace above this line \(too many statements above if\)`
	if b == 2 {
		panic(1)
	}

	_ = a
	_ = b
}

func fn2() {
	b := 2
	a := 1
	if b == 2 { // want `missing whitespace above this line \(no shared variables above if\)`
		panic(1)
	}

	_ = a
	_ = b
}

func fn3() {
	err := errors.New("error") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		panic(err)
	}
}

func fn4() {
	if a := 1; a != 2 {
		panic(a)
	}
	if a := 2; a != 2 { // want `missing whitespace above this line \(invalid statement above if\)`
		panic(a)
	}
}

func fn5(m any, k string) string {
	v := m.(map[string]string)
	if r, ok := m[k]; ok {
		return r
	}

	return k
}

func fn6() {
	if true {
		_ = 1
		_ = 2
		return // want `missing whitespace above this line \(too many lines above return\)`
	} else if false {
		_ = 1
		_ = 2
		return // want `missing whitespace above this line \(too many lines above return\)`
	} else {
		_ = 1
		_ = 2
		return // want `missing whitespace above this line \(too many lines above return\)`
	}
}

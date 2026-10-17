package testpkg

func bar() (string, bool) {
	return "test", true
}

func fn1() {
	var ok bool
	foo, ok := bar()

	_ = foo
	_ = ok
}

func fn2() {
	var ok bool
	var foo string // want `missing whitespace above this line \(never cuddle decl\)`

	foo, ok = bar()

	_ = foo
	_ = ok
}

func fn3() {
	var a int
	b := a + 1

	_ = b
}

func fn4() {
	var foo int
	bar := 2 // want `missing whitespace above this line \(invalid statement above assign\)`

	_ = foo
	_ = bar
}

func fn5() {
	var a = 1
	var b = 2 // want `missing whitespace above this line \(never cuddle decl\)`

	_, _ = a, b
}

package testpkg

func fn1() {
	var a = 1
	var b = 2 // want "missing whitespace decreases readability"

	const c = 3
	const d = 4 // want "missing whitespace decreases readability"

	e := 5
	var f = 6 // want "missing whitespace decreases readability"
	g := 7 // want "missing whitespace decreases readability"

	_ = a
	_ = b
	_ = c
	_ = d
	_ = e
	_ = f
	_ = g
}

func fn2() {
	var a = 1
	var b = a // want "missing whitespace decreases readability"
}

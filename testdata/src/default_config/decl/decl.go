package testpkg

func fn1() {
	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2

	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	const c = 3
	const d = 4

	// want +3 `missing whitespace above this line \(never cuddle decl\)`
	// want +3 `missing whitespace above this line \(invalid statement above assign\)`
	e := 5
	var f = 6
	g := 7

	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var h = 1
	var i = a

	_, _, _, _, _, _, _, _, _ = a, b, c, d, e, f, g, h, i
}

func fn2() {
	var x = func() int { // want +1 `unnecessary whitespace \(leading-whitespace\)`

		return 1
	}()

	_ = x
}

func fn3() {
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2
	var c = 3
	var d = 4

	_, _, _, _ = a, b, c, d
}

func fn4() {
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	var (
		a = 1
		b = 2
	)
	var (
		c = 3
		d = 4
	)

	_, _, _, _ = a, b, c, d
}

func fn5() {
	// want +7 `missing whitespace above this line \(never cuddle decl\)`
	// want +7 `missing whitespace above this line \(never cuddle decl\)`
	// want +9 `missing whitespace above this line \(never cuddle decl\)`
	var (
		a = 1
		b = 2
	)
	var x int
	var z = func() int {
		return 1
	}
	var (
		c = 3
		d = 4
	)

	_, _, _, _, _, _ = a, b, x, z, c, d
}

func fn6() {
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2
	var c = 3 // test
	var d = 4
	var e = 5

	_, _, _, _, _ = a, b, c, d, e
}

func fn7() {
	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2

	var c = 3

	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var d = 4
	var e = 5

	_, _, _, _, _ = a, b, c, d, e
}

func fn8() {
	// want +3 `missing whitespace above this line \(never cuddle decl\)`
	// Comment above
	var g = 7
	var h = 8
	// Comment after

	_, _ = g, h
}

func fn9() {
	type S1 struct {
		N int
	}
	type S2 struct { // want `missing whitespace above this line \(never cuddle decl\)`
		N int
	}
}

func fn10() {
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	a := 1
	b := 2
	var c int
	var d = "string"
	var e string = "string"
	f := 3 // want `missing whitespace above this line \(invalid statement above assign\)`

	_, _, _, _, _, _ = a, b, c, d, e, f
}

func fn11() {
	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var a int
	var b int
	if b > 0 {
		_ = 1
	}

	_ = a
}

func fn12() {
	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var a int
	var b int // Not grouped due to this comment
	if b > 0 {
		_ = 1
	}

	_ = a
}

func fn13() {
	var a int
	if a > 0 {
		_ = 1
	}
}

func fn14() {
	var (
		a = 1
		b = 2
		c = 3
	)
	if a+b > c {
		_ = 1
	}
}

func fn15() {
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	b := 2 // want `missing whitespace above this line \(invalid statement above assign\)`
	var c = 3
	var d = 4
	if a+b > c+d {
		_ = 1
	}
}

func fn16() {
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	b := 2 // want `missing whitespace above this line \(invalid statement above assign\)`
	var c = 3
	var d = 4 // Ungroupable
	if a+b > c+d {
		_ = 1
	}
}

func fn17() {
	// want +2 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2
	if true { // want `missing whitespace above this line \(no shared variables above if\)`
		_ = 1
	}

	_, _ = a, b
}

func fn18() {
	// want +5 `missing whitespace above this line \(never cuddle decl\)`
	var a = map[uint8]string{
		1: "foo",
		2: "bar",
	}
	var b = []struct {
		a string
		b int
		c bool
	}{
		{
			a: "string",
			b: "int",
			c: true,
		},
		{
			a: "string",
			b: "int",
			c: true,
		},
	}

	_, _ = a, b
}

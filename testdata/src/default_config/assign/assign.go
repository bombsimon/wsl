package testpkg

func strictAppend() {
	s := []string{}
	s = append(s, "a")
	s = append(s, "b")
	x := "c"
	s = append(s, x)
	y := "e"
	s = append(s, "d") // want "missing whitespace decreases readability"
	s = append(s, y)
}

func incDec() {
	x := 1
	x++
	x--
	y := x

	_ = y
}

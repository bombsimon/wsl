package testpkg

func fn1(s []string) {
	a := "a"
	s = append(s, a)

	x := 1
	s = append(s, "s") // want `missing whitespace above this line \(no shared variables above append\)`

	_ = x
}

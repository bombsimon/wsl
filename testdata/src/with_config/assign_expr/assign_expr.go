package testpkg

func assignAndCall(t1 SomeT) {
	t1.Fn()
	x := t1.Fn() // want `missing whitespace above this line \(invalid statement above assign\)`
	t1.Fn()
}

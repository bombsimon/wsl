package testpkg

func Fn(_ int) {}

func fn1() {
	one := 1
	two := 2
	three := 3

	a := 1
	if true {
		Fn(a)
	}

	b := 1
	if true { // want `missing whitespace above this line \(no shared variables above if\)`
		Fn(one)

		if false {
			Fn(two)

			if true {
				Fn(three)

				if true {
					Fn(b)
				}
			}
		}
	}
}

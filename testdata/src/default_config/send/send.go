package testpkg

func someFn(_ int) int {
	return 1
}

func fn1(ch chan int) {
	a := 1
	ch <- a

	b := 1
	ch <- someFn(b)

	c := 1
	ch <- 1 // want `missing whitespace above this line \(no shared variables above send\)`

	d := 2
	ch <- func() int { // want `missing whitespace above this line \(no shared variables above send\)`
		return c
	}()

	e := 2
	ch <- func() int {
		return e
	}()

	_ = d
}

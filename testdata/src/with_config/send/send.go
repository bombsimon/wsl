package testpkg

func fn1(ch chan int) {
	a := 1
	ch <- a

	b := 1
	ch <- someFn(b)

	c := 1
	ch <- 1 // want `missing whitespace decreases readability \(send\)`

	// TODO: Technically this is used first in block but there's no easy way to
	// figure out the first statement in the block so for now this is not valid.
	d := 2
	ch <- func() int { // want `missing whitespace decreases readability \(send\)`
		return d
	}()
}

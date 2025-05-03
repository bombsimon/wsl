package testpkg

func fn1(ch chan int) {
	a := 1
	ch <- 1 // want `missing whitespace decreases readability \(send\)`

	b := 2
	<-ch // want `missing whitespace decreases readability \(expr\)`
}

func fn2(ch chan int) {
	a := 1
	ch <- a

	b := 1
	b = <-ch
}

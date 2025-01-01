package testpkg

import "fmt"

func Go() {
	fooFunc := func() {}
	go fooFunc()

	barFunc := func() {}
	go fooFunc() // want "missing whitespace decreases readability"

	go func() {
		fmt.Println("hey")
	}()

	cuddled := true
	go func() { // want "missing whitespace decreases readability"
		fmt.Println("hey")
	}()

	argToGo := 1
	go takesArg(argToGo)

	notArgToGo := 1
	go takesArg(argToGo) // want "missing whitespace decreases readability"

	t1 := NewT()
	t2 := NewT()
	t3 := NewT()

	go t1()
	go t2()
	go t3()

	multiCuddle1 := NewT()
	multiCuddle2 := NewT() // want "missing whitespace decreases readability"
	go multiCuddle2()

	// TODO: Breaking change, this used to be on the first `go` stmt - now it's
	// on the line that should have a blank line above.
	t4 := NewT()
	t5 := NewT() // want "missing whitespace decreases readability"
	go t5()
	go t4()
}

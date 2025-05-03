package testpkg

import "fmt"

func NewT() func() {
	return func() {}
}

func Fn(_ int) {}

func Go() {
	fooFunc := func() {}
	go fooFunc()

	barFunc := func() {}
	go fooFunc() // want `missing whitespace decreases readability \(go\)`

	_ = barFunc

	go func() {
		fmt.Println("hey")
	}()

	cuddled := true
	go func() { // want `missing whitespace decreases readability \(go\)`
		fmt.Println("hey")
	}()

	_ = cuddled

	argToGo := 1
	go Fn(argToGo)

	notArgToGo := 1
	go Fn(argToGo) // want `missing whitespace decreases readability \(go\)`

	_ = notArgToGo

	t1 := NewT()
	t2 := NewT()
	t3 := NewT()

	go t1()
	go t2()
	go t3()

	multiCuddle1 := NewT()
	multiCuddle2 := NewT() // want `missing whitespace decreases readability \(go\)`
	go multiCuddle2()

	// TODO: Breaking change, this used to be on the first `go` stmt - now it's
	// on the line that should have a blank line above.
	t4 := NewT()
	t5 := NewT() // want `missing whitespace decreases readability \(go\)`
	go t5()
	go t4()

	_ = t1
	_ = t2
	_ = t3
	_ = t4
	_ = t5
	_ = multiCuddle1
	_ = multiCuddle2
}

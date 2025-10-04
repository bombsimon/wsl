package testpkg

import (
	"fmt"
	"os"
	"sync"
)

type T int

func Fn() T {
	return T(0)
}

func (*T) Close() {}

func fn() {
	a := Fn()
	b := Fn()

	defer a.Close()
	defer b.Close()

	a = Fn()
	defer a.Close()

	b = Fn()
	defer b.Close()

	a = Fn()
	defer a.Close()
	b = Fn() // want `missing whitespace above this line \(invalid statement above assign\)`
	defer b.Close()

	a = Fn()
	b = Fn()
	defer a.Close() // want `missing whitespace above this line \(no shared variables above defer\)`
	defer b.Close()

	m := sync.Mutex{}

	m.Lock()
	defer m.Unlock()

	c := true
	defer func(b bool) { // want `missing whitespace above this line \(no shared variables above defer\)`
		fmt.Printf("%v", b)
	}(false)

	_ = c
}

func fn2() {
	a := 1
	b := Fn() // want `missing whitespace above this line \(too many statements above defer\)`
	defer b.Close()

	_ = a
}

func fn3() {
	f, err := os.Open("x")
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

package testpkg

import (
	"errors"
	"fmt"
)

func errPairOverridesUnlimited() {
	err := errors.New("x")
	if err != nil {
		panic(err)
	}
}

func errPairOverridesUnlimitedWithExtra() {
	a := 1
	err := errors.New("x") // want `missing whitespace above this line \(too many statements above if\)`
	if err != nil {
		panic(err)
	}

	_ = a
}

func errPairUncuddledWithExtraUnlimited() {
	a := 1
	err := errors.New("x") // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		panic(err)
	}

	_ = a
}

func ifManyAllUsed() {
	a := 1
	b := 2
	c := 3
	d := 4
	if a+b+c+d > 0 {
		fmt.Println("ok")
	}
}

func ifEarlyNotUsed() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(variable not shared with if\)`
	c := 3
	d := 4
	if b+c+d > 0 {
		fmt.Println("ok")
	}

	_ = notUsed
}

func ifMiddleNotUsed() {
	a := 1
	notUsed := 2
	c := 3 // want `missing whitespace above this line \(variable not shared with if\)`
	if a+c > 0 {
		fmt.Println("ok")
	}

	_ = notUsed
}

func forManyAllUsed() {
	a := 1
	b := 2
	c := 3
	for i := a; i < b+c; i++ {
		fmt.Println(i)
	}
}

func forEarlyNotUsed() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(variable not shared with for\)`
	c := 3
	for i := b; i < c; i++ {
		fmt.Println(i)
	}

	_ = notUsed
}

func switchManyAllUsed() {
	a := 1
	b := 2
	c := 3
	switch {
	case a > b:
		fmt.Println("a")
	case c > a:
		fmt.Println("c")
	}
}

func goManyAllUsed() {
	a := 1
	b := 2
	c := 3
	go fmt.Println(a, b, c)
}

func goEarlyNotUsed() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(variable not shared with go\)`
	c := 3
	go fmt.Println(b, c)

	_ = notUsed
}

func deferManyAllUsed() {
	a := 1
	b := 2
	c := 3
	defer fmt.Println(a, b, c)
}

func deferEarlyNotUsed() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(variable not shared with defer\)`
	c := 3
	defer fmt.Println(b, c)

	_ = notUsed
}

func ifManyAllUsedInFirstBlock() {
	a := 1
	b := 2
	c := 3
	if true {
		fmt.Println(a, b, c)
	}
}

func ifEarlyNotUsedInFirstBlock() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(variable not shared with if\)`
	c := 3
	if true {
		fmt.Println(b, c)
	}

	_ = notUsed
}

func singleStillWorks() {
	a := 1
	if a > 0 {
		fmt.Println("ok")
	}
}

func noIntersectionStillErrors() {
	a := 1
	if true { // want `missing whitespace above this line \(no shared variables above if\)`
		fmt.Println("ok")
	}

	_ = a
}

func vars() {
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	// want +4 `missing whitespace above this line \(never cuddle decl\)`
	var a = 1
	var b = 2
	var c = 3
	var d = 4
	if a+b+c > d {
		fmt.Println("ok")
	}
}

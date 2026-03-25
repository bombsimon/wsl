package testpkg

import "fmt"

func ifVarsUsedInBody() {
	a := 1
	b := 2
	if true {
		fmt.Println(a, b)
	}
}

func ifOneVarNotUsedAnywhere() {
	notUsed := 1
	b := 2 // want `missing whitespace above this line \(too many statements above if\)`
	if true {
		fmt.Println(b)
	}

	_ = notUsed
}

func switchVarsUsedInBody() {
	a := 1
	b := 2
	c := 3
	switch {
	case true:
		fmt.Println(a, b, c)
	}
}

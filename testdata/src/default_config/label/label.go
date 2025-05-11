package testpkg

import "fmt"

func fn1() {
	for i := range make([]int, 2) {
		if i == 1 {
			goto END
		}

	END:
		_ = 1
	}
}

func fn2() {
	for i := range make([]int, 2) {
		if i == 1 {
			goto END
		}

	END:
	}
}

func fn3() {
L1:
	if true {
		_ = 1
	}
L2: // want `missing whitespace above this line \(never cuddle label\)`
	if true {
		_ = 1
	}

	_ = 1
L3: // want `missing whitespace above this line \(never cuddle label\)`
	_ = 1
}

func fn() {
LABEL:
	if true {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		break // want `missing whitespace above this line \(too many lines above branch\)`
	}
}

package testpkg

import "fmt"

func fn() {
LABEL:
	for range 1 {
		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			break // want `missing whitespace above this line \(too many lines above branch\)`
		}
	}

	goto LABEL
}

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

	goto L1

	goto L2

	goto L3
}

func fn4(err error) {
	if true { // want +1 `unnecessary whitespace \(leading-whitespace\)`

	LABEL:
		panic(err)

		goto LABEL
	}
}

func fn5() {
	a := 1
LABEL: // want `missing whitespace above this line \(never cuddle label\)`
	if true {
	}

	goto LABEL

	_ = a
}

func fn6() {
	i := 0
	i++
ADD: // want `missing whitespace above this line \(never cuddle label\)`
	i++
SUB: // want `missing whitespace above this line \(never cuddle label\)`
	i--

	goto ADD

	goto SUB
}

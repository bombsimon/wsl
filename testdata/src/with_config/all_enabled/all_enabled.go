package testpkg

import "fmt"

func e() (bool, error) {
	return false, nil
}

func blockFollowedByIf() {
	x, y := 1, 2
	if x > 0 {
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	if y > 0 { // want `missing whitespace above this line \(invalid statement above if\)`
		_ = 1
	}
}

func blockFollowedByFor() {
	x, y := 1, 2
	if x > 0 {
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	for i := 0; i < y; i++ { // want `missing whitespace above this line \(invalid statement above for\)`
		_ = i
	}
}

func blockFollowedBySwitch() {
	x, y := 1, 2
	if x > 0 {
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	switch y { // want `missing whitespace above this line \(invalid statement above switch\)`
	case 1:
		_ = 1
	}
}

func consecutiveBlocks() {
	if true {
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	if false { // want `missing whitespace above this line \(invalid statement above if\)`
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	for i := 0; i < 1; i++ { // want `missing whitespace above this line \(invalid statement above for\)`
		_ = i
	} // want `missing whitespace below this line \(newline-after-block\)`
	switch 1 { // want `missing whitespace above this line \(invalid statement above switch\)`
	case 1:
		_ = 1
	}
}

func caseAndBlockNewlines(n int) {
	switch n {
	case 1:
		_ = 1 // want `missing whitespace below this line \(case-trailing-newline\)`
	case 2:
		_ = 1
	} // want `missing whitespace below this line \(newline-after-block\)`
	fmt.Println("") // want `missing whitespace above this line \(invalid statement above expr\)`
}

func errCheckAndNewlineAfterBlock() {
	b, err := e() // want +1 `unnecessary whitespace \(err\)`

	if err != nil {
		fmt.Println("")
	} // want `missing whitespace below this line \(newline-after-block\)`
	fmt.Println(b) // want `missing whitespace above this line \(invalid statement above expr\)`
}

func multipleOnSameStatement() {
	a := 1
	b := 2 // want `missing whitespace above this line \(too many statements above if\)`
	if a > 0 {
		_ = b
	} // want `missing whitespace below this line \(newline-after-block\)`
	c := 3 // want `missing whitespace above this line \(invalid statement above assign\)`

	_ = c
}

func nestedBlocksAllNeedNewlines() {
	if true {
		if true {
			_ = 1
		} // want `missing whitespace below this line \(newline-after-block\)`
		if false { // want `missing whitespace above this line \(invalid statement above if\)`
			_ = 1
		} // want `missing whitespace below this line \(newline-after-block\)`
	} // want `missing whitespace below this line \(newline-after-block\)`
	_ = 1 // want `missing whitespace above this line \(invalid statement above assign\)`
}

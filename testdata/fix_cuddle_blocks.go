package testpkg

import "fmt"

func FailBlocks() {
	if false {
		return
	}
	_, err := maybeErr() // want "assignments should only be cuddled with other assignments"
	if err != nil {      // want "only one cuddle assignment allowed before if statement"
		return
	}

	if true {
		return
	}
	_, x := maybeErr() // want "assignments should only be cuddled with other assignments"
	if err != nil {    // want "only one cuddle assignment allowed before if statement"
		return
	}

	if x != nil {
	}
	switch 1 { // want "anonymous switch statements should never be cuddled"
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}
}

func maybeErr() (int, error) {
	return 0, nil
}

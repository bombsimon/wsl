package testpkg

import "fmt"

func FailBlocks() {
	if false {
		return
	}
	_, err := maybeErr()
	if err != nil {
		return
	}

	if true {
		return
	}
	_, x := maybeErr()
	if err != nil {
		return
	}

	if x != nil {
	}
	switch 1 {
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
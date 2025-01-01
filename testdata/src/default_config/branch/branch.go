package testpkg

import "fmt"

func fn1() {
	for range []int{} {
		if true {
			fmt.Println("")

			continue
		}

		if true {
			fmt.Println("")
			continue
		}

		if true {
			fmt.Println("")
			break
		}

		if false {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")

			continue
		}

		if false {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")

			break
		}
	}
}

func fn2() {
	for range []int{} {
		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			continue // want "missing whitespace decreases readability"
		}

		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			break // want "missing whitespace decreases readability"
		}
	}
}

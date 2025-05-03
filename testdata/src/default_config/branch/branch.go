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
			break
		}

		if true {
			fmt.Println("")
			continue
		}

		if true {
			fmt.Println("")
			fallthrough
		}

		if true {
			fmt.Println("")
			goto START
		}

		if false {
			fmt.Println("")
			fmt.Println("")
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

			fallthrough
		}

		if false {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")

			goto START
		}
	}
}

func fn2() {
	for range []int{} {
		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			break // want `missing whitespace decreases readability \(branch\)`
		}

		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			continue // want `missing whitespace decreases readability \(branch\)`
		}

		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fallthrough // want `missing whitespace decreases readability \(branch\)`
		}

		if true {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			goto START // want `missing whitespace decreases readability \(branch\)`
		}
	}
}

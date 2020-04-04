package testpkg

import "fmt"

func FailCaseTrailing() {
	switch 1 == 1 {
	case true:
		fmt.Println("true")
		// Error?
	case false:
		fmt.Println("false") // Error?!
	case false && true:
		fmt.Println("multitalented")
		// First Line
		// Second Line
	case false && false:
		fmt.Println("multitalented")
		/*
			This is
			Odd
		*/
	case true && false:
		fmt.Println("Compiles?")
	default:
		fmt.Println("idk..")
	}
}

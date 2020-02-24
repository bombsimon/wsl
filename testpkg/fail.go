package testpkg

import "fmt"

func ThisFuncFails() {
	one := 1
	two := 2     // Some comment
	if two > 3 { // More comment
		fmt.Println(one)
	}
}

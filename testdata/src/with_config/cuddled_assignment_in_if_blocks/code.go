package pkg

import (
	"fmt"
)

func ExampleFunction() {
	// This will cause an error when AllowCuddledAssignmentsInIfBlocks is false
	x := true
	if true {
		fmt.Println("didn't use cuddled variable")
	}

	_ = x // This assignment is outside the if statement's scope
}

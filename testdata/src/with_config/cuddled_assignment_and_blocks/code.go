package pkg

import (
	"fmt"
)

func IfFunction() {
	// This will cause an error when AllowCuddledAssignmentsInIfBlocks is false
	x := true
	if true {
		fmt.Println("didn't use cuddled variable")
	}

	_ = x // This assignment is outside the if statement's scope
}

func ForFunction() {
	// This will cause an error when AllowCuddledAssignmentsInIfBlocks is false
	x := true
	for {
		fmt.Println("didn't use cuddled variable")
	}

	_ = x // This assignment is outside the if statement's scope
}

package testpkg

import "fmt"

// Two sharing cuddled stmts is allowed (within max=2).
func ifTwoShareWithinLimit() {
	x := 1
	y := 2
	if x > y {
		fmt.Println("ok")
	}
}

// Three sharing stmts exceeds the limit: separate the whole group from the
// trigger (cuddle-group placement instead of splitting between vars).
func ifThreeShareSeparate() {
	a := 1
	b := 2
	c := 3
	if a+b+c > 0 { // want `missing whitespace above this line \(too many statements above if\)`
		fmt.Println("ok")
	}
}

// A non-sharing stmt in the chain still separates the whole group from the
// trigger, regardless of how many sharing stmts are within the limit.
func ifNonSharingAbove() {
	notUsed := 1
	a := 2
	b := 3
	if a+b > 0 { // want `missing whitespace above this line \(too many statements above if\)`
		fmt.Println("ok")
	}

	_ = notUsed
}

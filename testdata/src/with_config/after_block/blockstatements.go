package testpkg

import "fmt"

func ifStatementWithoutNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("next statement")
}

func ifStatementWithNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}

	fmt.Println("next statement")
}

func ifStatementAtEnd() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}
}

func ifElseStatement() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	}

	fmt.Println("next statement")
}

func ifElseStatementWithoutNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("next statement")
}

func ifElseStatementAtEnd() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	}
}

func ifElseIfStatement() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}

	fmt.Println("next statement")
}

func ifElseIfStatementWithoutNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("next statement")
}

func forLoopWithoutNewline() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after loop")
}

func forLoopWithNewline() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	fmt.Println("after loop")
}

func forLoopAtEnd() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
}

func rangeLoopWithoutNewline() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after loop")
}

func rangeLoopWithNewline() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("after loop")
}

func switchStatementWithoutNewline() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one") // want `missing whitespace below this line \(case-trailing-newline\)`
	case 2:
		fmt.Println("two") // want `missing whitespace below this line \(case-trailing-newline\)`
	default:
		fmt.Println("other")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after switch")
}

func switchStatementWithNewline() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")

	case 2:
		fmt.Println("two")

	default:
		fmt.Println("other")
	}

	fmt.Println("after switch")
}

func selectStatementWithoutNewline() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v) // want `missing whitespace below this line \(case-trailing-newline\)`
	default:
		fmt.Println("default")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after select")
}

func selectStatementWithNewline() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v)

	default:
		fmt.Println("default")
	}

	fmt.Println("after select")
}

func multipleStatementsWithMixedViolations() {
	x := 5

	if x > 0 {
		fmt.Println("positive")
	} // want `missing whitespace below this line \(after-block\)`
	for i := 0; i < x; i++ {
		fmt.Println(i)
	}

	switch x {
	case 5:
		fmt.Println("five")
	}
}

func nestedIfWithoutNewline() {
	x := 5
	y := 10
	if x > 0 {
		if y > 0 {
			fmt.Println("both positive")
		} // want `missing whitespace below this line \(after-block\)`
		fmt.Println("x positive")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after outer if")
}

func nestedIfWithNewline() {
	x := 5
	y := 10
	if x > 0 {
		if y > 0 {
			fmt.Println("both positive")
		}

		fmt.Println("x positive")
	}

	fmt.Println("after outer if")
}

func nestedForWithoutNewline() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Println(i, j)
		} // want `missing whitespace below this line \(after-block\)`
		fmt.Println("inner loop done")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after outer loop")
}

func nestedForWithNewline() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Println(i, j)
		}

		fmt.Println("inner loop done")
	}

	fmt.Println("after outer loop")
}

func nestedSwitchWithoutNewline() {
	x := 1
	y := 2
	switch x {
	case 1:
		switch y {
		case 2:
			fmt.Println("x=1, y=2")
		} // want `missing whitespace below this line \(after-block\)`
		fmt.Println("x=1") // want `missing whitespace below this line \(case-trailing-newline\)`
	case 2:
		fmt.Println("x=2")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after outer switch")
}

func nestedSwitchWithNewline() {
	x := 1
	y := 2
	switch x {
	case 1:
		switch y {
		case 2:
			fmt.Println("x=1, y=2")
		}

		fmt.Println("x=1")

	case 2:
		fmt.Println("x=2")
	}

	fmt.Println("after outer switch")
}

func complexNested() {
	for i := 0; i < 3; i++ {
		if i%2 == 0 {
			switch i {
			case 0:
				fmt.Println("zero") // want `missing whitespace below this line \(case-trailing-newline\)`
			case 2:
				fmt.Println("two")
			} // want `missing whitespace below this line \(after-block\)`
			fmt.Println("even")
		}

		fmt.Println("iteration", i)
	}

	fmt.Println("done")
}

func typeSwitchWithoutNewline() {
	a := any("hello")
	switch v := a.(type) {
	case string:
		fmt.Println("string:", v) // want `missing whitespace below this line \(case-trailing-newline\)`
	case int:
		fmt.Println("int:", v) // want `missing whitespace below this line \(case-trailing-newline\)`
	default:
		fmt.Println("unknown type")
	} // want `missing whitespace below this line \(after-block\)`
	fmt.Println("after type switch")
}

func typeSwitchWithNewline() {
	a := any("hello")
	switch v := a.(type) {
	case string:
		fmt.Println("string:", v)

	case int:
		fmt.Println("int:", v)

	default:
		fmt.Println("unknown type")
	}

	fmt.Println("after type switch")
}

func funcLiteralInvoked() {
	err := func() error {
		return nil
	}()
	if err != nil {
		panic(err)
	}
}

func funcLiteralAssigned() {
	someFn := func() {
		_ = 1
	}
	anotherVar := 1

	_ = someFn
	_ = anotherVar
}

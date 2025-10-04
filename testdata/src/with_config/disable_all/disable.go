package testpkg

import (
	"context"
	"fmt"
)

func Fn(func ()) {}

func assign() {
	var a = 1
	b := 1

	_ = a
	_ = b
}

func branch() {
	for range make([]int, 2) {
		_ = 1
		_ = 1
		_ = 1
		continue
	}

	for range make([]int, 2) {
		_ = 1
		_ = 1
		_ = 1
		break
	}
}

func decl() {
	var a = 1
	var b = 1
	var c = b

	_ = a
	_ = b
	_ = c
}

func deferFunc() {
	d := 1
	defer func() {}()

	_ = d
}

func expr() {
	a := 1
	b := 2
	fmt.Println("")
	c := 3
	d := 4

	_, _, _, _ = a, b, c, d
}

func forLoop() {
	a := 1
	b := 2
	for i := 0; i < b; i++ {
		panic(1)
	}

	for i := 0; i < 1; i++ {
		panic("")
	}
	for i := 0; i < 1; i++ {
		panic("")
	}

	_ = a
	_ = b
}

func goRoutine() {
	f1 := func() {}

	f2 := func() {}
	go f1()

	f3 := func() {}
	go Fn(f1)

	_, _, _ = f1, f2, f3
}

func ifStmt() {
	a := 1
	b := 2
	if b == 2 {
		panic(1)
	}

	b = 2
	a = 1
	if b == 2 {
		panic(1)
	}

	_ = a
	_ = b
}

func label() {
L1:
	if true {
		_ = 1
	}
L2:
	if true {
		_ = 1
	}

	goto L1

	goto L2
}

func rangeLoop() {
	a := []int{}
	b := []int{}
	for range b {
		panic(1)
	}

	_ = a
	_ = b
}

func returnStmt() int {
	if true {
		_ = 1
		_ = 2
		return 1
	}
	return 1
}

func selectStmt(ctx context.Context, ch1 chan struct{}) {
	x := 1
	select {
	case <-ctx.Done():
		_ = 1
	case <-ch1:
		_ = 1
	}

	_ = x
}

func send(ch chan int) {
	a := 1
	ch <- 1

	_ = a
}

func switchStmt() {
	a := 1
	b := 2
	switch b {
	case 1:
	case 2:
	case 3:
		panic(a)
	}
}

func typeSwitch() {
	var a any = 1

	b := 1
	switch a.(type) {
	case int:
	case string:
	}

	_ = b
}

func whitespace() {

	if true {

		_ = 1

	}

}

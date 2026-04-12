package testpkg

import "fmt"

type item struct {
	key string
}

func (i *item) setKey(k string) {
	i.key = k
}

// ifWithIntersection: expr stmt sharing a variable with if — allowed.
func ifWithIntersection(i *item) {
	i.setKey("hello")
	if i.key == "hello" {
		fmt.Println("match")
	}
}

// ifWithoutIntersection: expr stmt not sharing a variable with if — still reported.
func ifWithoutIntersection(a int, b int) {
	fmt.Println(b)
	if a > 0 { // want `missing whitespace above this line \(no shared variables above if\)`
		_ = a
	}
}

// forWithIntersection: expr stmt sharing a variable with for — allowed.
func forWithIntersection(items []item, key string) {
	items[0].setKey(key)
	for _, it := range items {
		fmt.Println(it.key)
	}
}

// goWithIntersection: expr stmt sharing a variable with go — allowed.
func goWithIntersection(i *item) {
	i.setKey("async")
	go fmt.Println(i.key)
}

// deferWithIntersection: expr stmt sharing a variable with defer — allowed.
func deferWithIntersection(i *item) {
	i.setKey("deferred")
	defer fmt.Println(i.key)
}

func nestedExprBeforeInnerIf(i *item, key string, check bool) {
	if check {
		i.setKey(key)
		if i.key == key {
			fmt.Println("match")
		}
	}
}

func exprBetweenAssignAndIf(a, b int) {
	someArg := "hello"
	fmt.Println(someArg)
	if a < b {
		someArg = "world"
		fmt.Println(someArg)
	}
}

func mixingExpr(b, d int) {
	a := 1
	fmt.Println(b)
	c := 2
	fmt.Println(d)
	if a+b+c+d > 0 {
		fmt.Println("mix")
	}
}

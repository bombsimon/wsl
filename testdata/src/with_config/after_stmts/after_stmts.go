package testpkg

import (
	"fmt"
	"os"
	"sync"
)

// --- after-defer ---

func simpleDeferMissing() {
	defer fmt.Println("cleanup") // want `missing whitespace below this line \(after-defer\)`
	x := 5

	fmt.Println(x)
}

func simpleDeferOK() {
	defer fmt.Println("cleanup")

	x := 5

	fmt.Println(x)
}

func multipleDefersMissing() {
	defer fmt.Println("a")
	defer fmt.Println("b") // want `missing whitespace below this line \(after-defer\)`
	x := 5

	fmt.Println(x)
}

func multipleDefersOK() {
	defer fmt.Println("a")
	defer fmt.Println("b")

	x := 5

	fmt.Println(x)
}

func deferAfterErrCheckMissing() {
	f, err := os.Open("example.txt")
	if err != nil {
		_ = err
	}
	defer f.Close() // want `missing whitespace below this line \(after-defer\)`
	data := []byte("test")
	_ = f
	_ = data
}

func deferLastInFunc() {
	x := 1

	fmt.Println(x)

	defer fmt.Println("cleanup")
}

// --- after-go ---

func simpleGoMissing() {
	go fmt.Println("work") // want `missing whitespace below this line \(after-go\)`
	x := 5

	fmt.Println(x)
}

func simpleGoOK() {
	go fmt.Println("work")

	x := 5

	fmt.Println(x)
}

func multipleGoMissing() {
	go fmt.Println("a")
	go fmt.Println("b") // want `missing whitespace below this line \(after-go\)`
	x := 5

	fmt.Println(x)
}

func multipleGoOK() {
	go fmt.Println("a")
	go fmt.Println("b")

	x := 5

	fmt.Println(x)
}

func goAnonFuncMissing() {
	go func() {
		fmt.Println("work")
	}() // want `missing whitespace below this line \(after-go\)`
	fmt.Println("next")
}

// --- after-decl ---

func simpleVarMissing() {
	var x int // want `missing whitespace below this line \(after-decl\)`
	x = 1

	fmt.Println(x)
}

func simpleVarOK() {
	var x int

	x = 1

	fmt.Println(x)
}

func cuddledVarsMissing() {
	var a int
	var b int // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(a, b)
}

func cuddledVarsOK() {
	var a int
	var b int

	fmt.Println(a, b)
}

func varBlockMissing() {
	var (
		a = 1
		b = 2
	) // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(a, b)
}

func varBlockOK() {
	var (
		a = 1
		b = 2
	)

	fmt.Println(a, b)
}

func constMissing() {
	const c = 1 // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(c)
}

func declAsOnlyStmt() {
	var _ = 1
}

// --- after-expr ---

func simpleExprMissing() {
	fmt.Println("hello") // want `missing whitespace below this line \(after-expr\)`
	x := 5

	fmt.Println(x)
}

func simpleExprOK() {
	fmt.Println("hello")

	x := 5

	fmt.Println(x)
}

func multipleExprsMissing() {
	fmt.Println("a")
	fmt.Println("b") // want `missing whitespace below this line \(after-expr\)`
	x := 5

	fmt.Println(x)
}

func multipleExprsOK() {
	fmt.Println("a")
	fmt.Println("b")

	x := 5

	fmt.Println(x)
}

// --- combined interactions ---

func deferFollowedByExpr() {
	defer fmt.Println("cleanup") // want `missing whitespace below this line \(after-defer\)`
	fmt.Println("next")          // want `missing whitespace below this line \(after-expr\)`
	x := 1
	_ = x
}

func goFollowedByDefer() {
	go fmt.Println("work")       // want `missing whitespace below this line \(after-go\)`
	defer fmt.Println("cleanup") // want `missing whitespace below this line \(after-defer\)`
	x := 1
	_ = x
}

func exprFollowedByDeferWithIntersection() {
	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()

	fmt.Println("protected")
}

func exprFollowedByDeferNoIntersection() {
	var mu sync.Mutex

	mu.Lock() // want `missing whitespace below this line \(after-expr\)`
	defer fmt.Println("cleanup")

	fmt.Println("world")
}

func declFollowedByExpr() {
	var x = 1 // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(x)
}

func mixedDeclKindsMissing() {
	var a = 1 // want `missing whitespace below this line \(after-decl\)`
	type B struct {
		B string
	} // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(a) // want `missing whitespace below this line \(after-expr\)`
	_ = B{}
}

func mixedDeclKindsOK() {
	var a = 1

	type B struct {
		B string
	}

	fmt.Println(a)

	_ = B{}
}

func consecutiveTypesMissing() {
	type A struct {
		A string
	}
	type B struct {
		B string
	} // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(A{}, B{})
}

func consecutiveTypesOK() {
	type A struct {
		A string
	}
	type B struct {
		B string
	}

	fmt.Println(A{}, B{})
}

func varThenConstMissing() {
	var a = 1 // want `missing whitespace below this line \(after-decl\)`
	const b = 2 // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(a, b)
}

func varThenConstOK() {
	var a = 1

	const b = 2

	fmt.Println(a, b)
}

func constThenVarMissing() {
	const a = 1 // want `missing whitespace below this line \(after-decl\)`
	var b = 2 // want `missing whitespace below this line \(after-decl\)`
	fmt.Println(a, b)
}

func sameKindCuddlesOK() {
	var a = 1
	var b = 2

	const c = 3
	const d = 4

	fmt.Println(a, b, c, d)
}

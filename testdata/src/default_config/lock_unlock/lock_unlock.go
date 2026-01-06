package testpkg

import (
	"log"
	"sync"
)

type S struct {
	mu sync.Mutex
}

func mutex(mu sync.Mutex, items []int) {
	mu.Lock()
	for _, i := range items {
		_ = i
	}
	mu.Unlock()
}

func mutexOnStruct(s S, items []int) {
	s.mu.RWLock()
	for _, i := range items {
		_ = i
	}
	s.mu.RWUnlock()
}

func mutexWithoutBlockMultipleVars(mu sync.Mutex) {
	mu.Lock()
	x := 1
	y := 2
	mu.Unlock()

	_, _ = x, y
}

func nonLockExpr(mu log.Logger) {
	mu.Print("before")
	x := 1            // want `missing whitespace above this line \(invalid statement above assign\)`
	mu.Print("after") // want `missing whitespace above this line \(no shared variables above expr\)`

	_ = x
}

func otherAssignments(mu sync.Mutex, items []int) {
	foo := 1
	mu.Lock() // want `missing whitespace above this line \(no shared variables above expr\)`
	for _, i := range items {
		_ = i
	}
	mu.Unlock()
	bar := 1 // want `missing whitespace above this line \(invalid statement above assign\)`

	_, _ = foo, bar
}

func lockInBlock(mu sync.Mutex, items []int) {
	func() {
		mu.Lock()
	}()
	for _, i := range items { // want `missing whitespace above this line \(invalid statement above range\)`
		_ = i
	}
	func() { // want `missing whitespace above this line \(invalid statement above expr\)`
		mu.Unlock()
	}()
}

func withDefer(mu sync.Mutex, items []int) {
	mu.Lock()
	defer mu.Unlock()
	for _, i := range items { // want `missing whitespace above this line \(invalid statement above range\)`
		_ = i
	}
}

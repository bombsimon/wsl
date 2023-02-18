package pkg

import "fmt"

func AllowAssignAndCallCuddle() {
	p.token(':')
	d.StartBit = p.uint() // want "assignments should only be cuddled with other assignments"
	p.token('|')          // want "only cuddled expressions if assigning variable or using from line above"
	d.Size = p.uint()     // want "assignments should only be cuddled with other assignments"
}

func AllowMultiLineAssignCuddle() {
	err := SomeFunc(
		"taking multiple values",
		"suddendly not a single line",
	)
	if err != nil { // want "if statements should only be cuddled with assignments"
		fmt.Println("denied")
	}
}

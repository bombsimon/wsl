package pkg

func AllowAssignAndCallCuddle() {
	p.token(':')
	d.StartBit = p.uint() // want "assignments should only be cuddled with other assignments"
	p.token('|')          // want "only cuddled expressions if assigning variable or using from line above"
	d.Size = p.uint()     // want "assignments should only be cuddled with other assignments"
}

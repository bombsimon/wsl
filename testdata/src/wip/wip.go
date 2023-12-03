package testpkg

func fn() {
	_ = 1

	// Comment before 1

	// Comment before 2
	var a = 1 // want "declarations should never be cuddled"
	var b = 2
	// Comment after 1

	// Comment after 2

	// One
	var a = 1 // want "declarations should never be cuddled"
	var b = 2
	// Two

	_ = 1
}

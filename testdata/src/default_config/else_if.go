func fn() {
	if true { // want "block should not start with a whitespace"

		fmt.Println("a")
	}

	if true {
		fmt.Println("a")
	} else if false { // want "block should not start with a whitespace"

		fmt.Println("b")
	} else { // want "block should not start with a whitespace"

		fmt.Println("c")
		fmt.Println("c")
		fmt.Println("c")
		return // want "return statements should not be cuddled if block has more than two lines"

	} // want "block should not end with a whitespace"
}

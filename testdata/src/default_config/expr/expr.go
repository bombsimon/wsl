package testpkg

import "fmt"

func fn() {
	a := 1
	b := 2
	c := 3
	a = 4
	b = 3

	_ = a
	_ = b
	_ = c
}

func fn2() {
	a := 1
	b := 2
	fmt.Println("") // want `missing whitespace above this line \(no shared variables above expr\)`
	c := 3          // want `missing whitespace above this line \(invalid statement above assign\)`
	d := 4

	_ = a
	_ = b
	_ = c
	_ = d
}

func fn3() {
	a := 1
	b := 2
	fmt.Println(b)
	c := 3 // want `missing whitespace above this line \(invalid statement above assign\)`
	d := 4

	_ = a
	_ = b
	_ = c
	_ = d
}

func fn4() {
	b := 2
	<-ch // want `missing whitespace above this line \(no shared variables above expr\)`
}

func fn5() {
	s := []string{
		func() string { // want +1 `unnecessary whitespace \(leading-whitespace\)`

			"a"
		},
		"b",
	}
}

func issue153(b bool) {
	for i := 0; i < 4; i++ {
		fmt.Println()
	}
	Up(3) // want `missing whitespace above this line \(invalid statement above expr\)`

	if autoheight != 3 {
		fmt.Printf("height should be 3 but is %d", autoheight)
	}
	Down(3) // want `missing whitespace above this line \(invalid statement above expr\)`

	if b {
		fmt.Println()
	}
	Up(b) // want `missing whitespace above this line \(invalid statement above expr\)`
}

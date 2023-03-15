package pkg

func Fn() {
	a = "a"
	b = "b"
	c = "c"

	a := "a"
	b := "b"
	c := "c"

	a := "a" // want "short declaration should cuddle only with other short declarations"
	b = "b"
	c := "c" // want "short declaration should cuddle only with other short declarations"

	a := "a" // want "short declaration should cuddle only with other short declarations"
	b = "b"  // Hey
	c := "c" // want "short declaration should cuddle only with other short declarations"
}

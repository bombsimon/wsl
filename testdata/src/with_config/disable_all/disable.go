package testpkg

func assign() {
	var a = 1
	b := 1

	_ = a
	_ = b
}

func branch() {
	for i := range(make(int[], 2)) {
		_ = 1
		_ = 1
		_ = 1
		continue
	}

	for i := range(make(int[], 2)) {
		_ = 1
		_ = 1
		_ = 1
		break
	}
}

func decl() {
	var a = 1
	var b = 1

	_ = a
	_ = b
}

func defer() {
	d := 1
	defer func() {}()
}

func expr() {
	a := 1
	b := 2
	fmt.Println("")
	c := 3
	d := 4
}

func for() {
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

func go() {
	f1 := func() {}

	f2 := func() {}
	go f1()

	f3 := func() {}
	go Fn(f1)
}

func if() {
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
}

func range() {
	a := []int{}
	b := []int{}
	for range b {
		panic(1)
	}

	_ = a
	_ = b
}

func return() {
	if true {
		_ = 1
		_ = 2
		return 1
	}
	return 1
}

func select(ctx context.Context, ch1 chan struct{}) {
	x := 1
	select {
	case ctx.Done():
		_ = 1
	case <-ch1:
		_ = 1
	}

	_ = x
}

func send(ch chan int) {
	a := 1
	ch <- 1

	b := 2
	<-ch
}

func switch() {
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

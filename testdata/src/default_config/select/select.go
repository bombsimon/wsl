package testpkg

import "context"

func fn1(ctx context.Context, ch1 chan struct{}) {
	select {
	case ctx.Done():
		_ = 1
	case <-ch1:
		_ = 1
	default:
		_ = 1
	}
}

func fn1(ctx context.Context, ch1 chan struct{}) {
	x := 1
	select { // want `missing whitespace above this line \(no shared variables above select\)`
	case ctx.Done():
		_ = 1
	case <-ch1:
		_ = 1
	default:
		_ = 1
	}

	_ = x
}

func fn2(ch1 chan struct{}) {
	ch := make(chan struct{})
	select {
	case <-ch1:
		_ = 1
	case <-ch:
		_ = 1
	}

	ctx := getCtx()
	select {
	case <-ch1:
		_ = 1
	case <-ctx.Done():
		_ = 1
	}

	ch2 := make(chan struct{})
	ch3 := make(chan struct{}) // want `missing whitespace above this line \(too many statements above select\)`
	select {
	case <-ch1:
		_ = 1
	case <-ch2:
		_ = 1
	case <-ch3:
		_ = 1
	}
}

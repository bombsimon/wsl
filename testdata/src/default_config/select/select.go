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
	select { // want "missing whitespace decreases readability"
	case ctx.Done():
		_ = 1
	case <-ch1:
		_ = 1
	default:
		_ = 1
	}

	_ = x
}

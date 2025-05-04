package testdata

func fn1(s map[string]bool) {
	// TODO: Regression v1
	timeout := 10
	switch { // want `missing whitespace decreases readability`
	case s["bad timeout"]:
		timeout = 100
	case s["zero timeout"]:
		timeout = 0
	}
}

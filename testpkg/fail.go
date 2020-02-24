package testpkg

func ThisFuncFails() {
	if false {
		return
	}
	_, err := maybeErr()
	if err != nil {
		return
	}

	if true {
		return
	}
	_, x := maybeErr()
	if err != nil {
		return
	}

	if x != nil {
	}
}

func maybeErr() (int, error) {
	return 0, nil
}

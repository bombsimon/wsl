package pkg

func Fn() {
	err := ErrorProducingFunc()

	if err != nil { // want "if statements that check an error must be cuddled with the statement that assigned the error"
		return err
	}
}

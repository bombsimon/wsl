package testpkg

func Advanced() {
	var foo = 1
	var bar = 2    // want "declarations should never be cuddled"
	var biz int // want "declarations should never be cuddled"

	x := []string{}
	x = append(x, "literal")
	notUsed := "just assigning, don't mind me"
	x = append(x, "not z..") // want "append only allowed to cuddle with appended value"
	useMe := "right away"
	alsoNotUsed := ":("
	x = append(x, "just noise") // want "append only allowed to cuddle with appended value"
	x = append(x, useMe)

	_ = foo
	_ = bar
	_ = biz
	_ = notUsed
	_ = alsoNotUsed
}

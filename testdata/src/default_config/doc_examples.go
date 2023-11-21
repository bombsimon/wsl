package pkg

import (
	"fmt"
	"os"
)

func AnonyMousStatement() {
	timeout := 10
	switch { // want "anonymous switch statements should never be cuddled"
	case true:
		timeout = 100
	case false:
		timeout = 0
	}
}

func Append() {
	x := []string{}
	x = append(x, "literal")
	notUsed := "just assigning, don't mind me"
	x = append(x, "not z..") // want "append only allowed to cuddle with appended value"
	useMe := "right away"
	alsoNotUsed := ":("
	x = append(x, "just noise") // want "append only allowed to cuddle with appended value"
	x = append(x, useMe)
}

func Assignment() {
	if x == 1 {
		x = 0
	}
	z := x + 2 // want "assignments should only be cuddled with other assignments"

	fmt.Println("x")
	y := "x" // want "assignments should only be cuddled with other assignments"
}

func TrailingWhitespaceOrComment() {
	return fmt.Sprintf("x")
	// REMEMBER: add mux function later.

} // want "block should not end with a whitespace"

// Default is for this turned off so no diagnostic.
func CaseTrailingWhitespace() {
	switch n {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	}
}

func StartingWhitespace() { // want "block should not start with a whitespace"


	return
}

func BranchStatement() {
	for i := range make([]int, 5) {
		if i > 2 {
			sendToOne(i)
			sendToSecond(i)
			continue // want "branch statements should not be cuddled if block has more than two lines"
		}

		if statement == "is short" {
			sendToOne(i)
			break
		}
	}
}

func Declarations() {
	var eol string
	var i int // want "declarations should never be cuddled"
}

func Defer() {
	x := 4
	defer notX() // want "defer statements should only be cuddled with expressions on same variable"

	thisIsX := func() {}
	defer thisIsX()
}

func ExpressionAndBlocks() {
	t, err := someFunc()
	if err != nil {
		// handle error
		return
	}
	fmt.Println(t) // want "expressions should not be cuddled with blocks"
}

func ExpressionAndDeclarationOrReturn() {
	var i int
	run() // want "expressions should not be cuddled with declarations or returns"

	return i
	run() // want "expressions should not be cuddled with declarations or returns"
}

func ForWithoutCondition() {
	x := "hey"
	for { // want "for statement without condition should never be cuddled"
		fmt.Printf("x")
	}
}

func For() {
	x := y + 2
	for i := 0; i < y+10; i++ { // want "for statements should only be cuddled with assignments used in the iteration"
		fmt.Println(i)
	}

	list := getList()
	for i, v := range anotherList { // want "ranges should only be cuddled with assignments used in the iteration"
		fmt.Println(i)
	}
}

func Go() {
	thisIsFunc := func() {}
	go thisIsNOTFunc() // want "go statements can only invoke functions assigned on line above"

	thisIsFunc := func() {}
	go thisIsFunc()
}

func IfAssign() {
	one := 1
	if 1 == 1 { // want "if statements should only be cuddled with assignments used in the if statement itself"
		fmt.Println("duuh")
	}
	if 1 != one { // want "if statements should only be cuddled with assignments"
		fmt.Println("i don't believe you!")
	}

	two := 2
	if two != one {
		fmt.Println("i don't believe you!")
	}
}

func ExpressionAndAssignment() {
	notInExpression := GetIt()
	call.SomethingElse() // want "only cuddled expressions if assigning variable or using from line above"

	thisIsBetter := GetIt()
	thisIsBetter.CallIt()
}

func OnlyOneDefer() {
	assignmentOne := Something()
	assignmentTwo := AnotherThing()
	defer assignmentTwo() // want "only one cuddle assignment allowed before defer statement"

	f1, err := os.Open("/path/to/f1.txt")
	if err != nil {
		return -1
	}
	defer f1.Close()

	m.Lock()
	defer m.Unlock()
}

func OnlyOneFor() {
	i := 0
	a := 0
	for i = 0; i < 5; i++ { // want "only one cuddle assignment allowed before for statement"
		fmt.Println("x")
	}

	a := 0
	i := 0
	for i = 0; i < 5; i++ { // want "only one cuddle assignment allowed before for statement"
		fmt.Println("x")
	}
}

func OnlyOneGo() {
	assignOne := SomeOne()
	assignTwo := SomeTwo()
	go assignTwo() // want "only one cuddle assignment allowed before go statement"
}

func OnlyOneIf() {
	lengthA := len(sliceA)
	lengthB := len(sliceB)
	if lengthA != lengthB { // want "only one cuddle assignment allowed before if statement"
		fmt.Println("not equal")
	}
}

func OnlyOneRange() {
	listA := GetA()
	listB := GetB()
	for _, v := range listB { // want "only one cuddle assignment allowed before range statement"
		fmt.Println(v)
	}
}

func OnlyOneSwitch() {
	assignOne := SomeOne()
	assignTwo := SomeTwo()
	switch assignTwo { // want "only one cuddle assignment allowed before switch statement"
	case 1:
		fmt.Println("one")
	default:
		fmt.Println("NOT one")
	}
}

func OnlyOneTypeSwitch() {
	someT := GetT()
	anotherT := GetT()
	switch v := anotherT.(type) { // want "only one cuddle assignment allowed before type switch statement"
	case int:
		fmt.Println("was int")
	default:
		fmt.Println("was not int")
	}
}

func RangeWithAssignment() {
	y := []string{"a", "b", "c"}

	x := 5
	for _, v := range y { // want  "ranges should only be cuddled with assignments used in the iteration"
		fmt.Println(v)
	}

	x := 5

	y := []string{"a", "b", "c"}
	for _, v := range y {
		fmt.Println(v)
	}
}

func Return() {
	_ = func() {
		y += 15
		y += 15
		return fmt.Sprintf("%s", y) // want "return statements should not be cuddled if block has more than two lines"
	}

	_ = func() {
		y += 15
		return fmt.Sprintf("%s", y)
	}
}

func Switch() {
	notInSwitch := ""
	switch someThingElse { // want "switch statements should only be cuddled with variables switched"
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("one")
	case 3:
		fmt.Println("three")
	}

	inSwitch := 2
	switch inSwitch {
	case 1:
		fmt.Println("better")
	}
}

func TypeSwitch() {
	notInSwitch := GetSomeType()
	switch someThingElse.(type) { // want "type switch statements should only be cuddled with variables switched"
	case int:
		fmt.Println("int")
	default:
		fmt.Println("not int")
	}

	inSwitch := GetSomeType()
	switch inSwitch.(type) {
	case int:
		fmt.Println("better")
	}
}

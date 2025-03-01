package testpkg

import (
	"fmt"
	"strconv"
)

func ForCuddleAssignmentWholeBlock() {
	x := 1
	y := 2

	var numbers []int
	for i := 0; i < 10; i++ {
		if x == y {
			numbers = append(numbers, i)
		}
	}

	var numbers2 []int
	for i := range 10 {
		if x == y {
			numbers2 = append(numbers2, i)
		}
	}

	var numbers3 []int
	for {
		if x == y {
			numbers3 = append(numbers3, 1)
		}
	}

	var numbers4 []int
	for {
		numbers4 = append(numbers4, 1)
	}

	environment = make(map[string]string)
	for _, env := range req.GetConfig().GetEnvs() {
		switch env.GetKey() {
		case "user-data":
			cloudInitUserData = env.GetValue()
		default:
			environment[env.GetKey()] = env.GetValue()
		}
	}
}

func IfCuddleAssignmentWholeBlock() {
	x := 1
	y := 2

	counter := 0
	if somethingTrue {
		checker := getAChecker()
		if !checker {
			return
		}

		counter++ // Cuddled variable used in block, but not as first statement
	}

	var number2 []int
	if x == y {
		fmt.Println("a")
	} else {
		if x > y {
			number2 = append(number2, i)
		}
	}

	var number3 []int
	if x == y {
		fmt.Println("a")
	} else if x > y {
		if x == y {
			number3 = append(number3, i)
		}
	}

	var number4 []int
	if x == y {
		if x == y {
			number4 = append(number4, i)
		}
	} else if x > y {
		if x == y {
			number4 = append(number4, i)
		}
	} else {
		if x > y {
			number4 = append(number4, i)
		}
	}
}

func SwitchCuddleAssignmentWholeBlock() {
	var id string
	var b bool // want "declarations should never be cuddled"
	switch b { // want "only one cuddle assignment allowed before switch statement"
	case true:
		id = "a"
	case false:
		id = "b"
	}

	var b bool
	var id string // want "declarations should never be cuddled"
	switch b {    // want "only one cuddle assignment allowed before switch statement"
	case true:
		id = "a"
	case false:
		id = "b"
	}

	var id2 string
	switch i := objectID.(type) {
	case int:
		if true {
			id2 = strconv.Itoa(i)
		}
	case uint32:
		if true {
			id2 = strconv.Itoa(int(i))
		}
	case string:
		if true {
			id2 = i
		}
	}

	var id3 string
	switch {
	case int:
		if true {
			id3 = strconv.Itoa(i)
		}
	case uint32:
		if true {
			id3 = strconv.Itoa(int(i))
		}
	case string:
		if true {
			id3 = i
		}
	}
}

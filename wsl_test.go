package wsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericHandling(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "a file must start with package declaration",
			code: []byte(`func a() {
				fmt.Println("this is function a")
			}`),
			expectedErrorStrings: []string{"invalid syntax, file cannot be linted"},
		},
		{
			description: "comments are ignored, ignore this crappy code",
			code: []byte(`package main

			/*
			Multi line comments ignore errors
			if true {

				fmt.Println("Well hello blankies!")

			}
			*/
			func main() {
				// Also signel comments are ignored
				// foo := false
				// if true {
				//
				//	foo = true
				//
				// }
			}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r := ProcessFile("unit-test", tc.code)

			require.Len(t, r, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, r[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestShouldRemoveEmptyLines(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "if statements cannot start or end with newlines",
			code: []byte(`package main
			func main() {
				if true {

					fmt.Println("this violates the rules")

				}
			}`),
			expectedErrorStrings: []string{"block should not start with a whitespace", "block should not end with a whitespace (or comment)"},
		},
		{
			description: "whitespaces parsed correct even with comments (whitespace found)",
			code: []byte(`package main
			func main() {
				if true {
					// leading comment, still too much whitespace

					fmt.Println("this violates the rules")

					// trailing comment, still too much whitespace
				}
			}`),
			expectedErrorStrings: []string{"block should not start with a whitespace", "block should not end with a whitespace (or comment)"},
		},
		{
			description: "whitespaces parsed correct even with comments (whitespace NOT found)",
			code: []byte(`package main

			func main() {
				if true {
					// leading comment, still too much whitespace
					// more comments
					fmt.Println("this violates the rules")
				}
			}`),
		},
		{
			description: "whitespaces parsed correctly in case blocks",
			code: []byte(`package main
			func main() {
				switch true {
				case 1:

					fmt.Println("this is too much...")

				case 2:
					fmt.Println("this is too much...")

				}
			}`),
			expectedErrorStrings: []string{
				"block should not end with a whitespace (or comment)",
				"case block should not start with a whitespace",
				"case block should not end with a whitespace",
			},
		},
		{
			description: "whitespaces parsed correctly in case blocks with comments (whitespace NOT found)",
			code: []byte(`package main
			func main() {
				switch true {
				case 1:
					// Some leading comment here
					fmt.Println("this is too much...")
				case 2:
					fmt.Println("this is too much...")
				}
			}`),
		},
		{
			description: "declarations may never be cuddled",
			code: []byte(`package main

			func main() {
				var b = true

				if b {
					// b is true here
				}
				var a = false
			}`),
			expectedErrorStrings: []string{"declarations should never be cuddled"},
		},
		{
			description: "nested if statements should not be seen as cuddled",
			code: []byte(`package main

			func main() {
				if true {
					if false {
						if true && false {
							// Whelp?!
						}
					}
				}
			}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r := ProcessFile("unit-test", tc.code)

			require.Len(t, r, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, r[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestShouldAddEmptyLines(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "ok to cuddled if statements directly after assigments multiple values",
			code: []byte(`package main

			func main() {
				v, err := strconv.Atoi("1")
				if err != nil {
					fmt.Println(v)
				}
			}`),
		},
		{
			description: "ok to cuddled if statements directly after assigments single value",
			code: []byte(`package main

			func main() {
				a := true
				if a {
					fmt.Println("I'm OK")
				}
			}`),
		},
		{
			description: "ok to cuddled if statements directly after assigments single value, negations",
			code: []byte(`package main

			func main() {
				a, b := true, false
				if !a {
					fmt.Println("I'm OK")
				}
			}`),
		},
		{
			description: "cannot cuddled if assigments not used",
			code: []byte(`package main

			func main() {
				err := ProduceError()

				a, b := true, false
				if err != nil {
					fmt.Println("I'm OK")
				}
			}`),
			expectedErrorStrings: []string{"if statements should only be cuddled with assigments used in the if statement itself"},
		},
		{
			description: "cannot cuddled with other things than assigments",
			code: []byte(`package main

			func main() {
				if true {
					fmt.Println("I'm OK")
				}
				if false {
					fmt.Println("I'm OK")
				}
			}`),
			expectedErrorStrings: []string{"if statements should only be cuddled with assigments"},
		},
		{
			description: "if statement with multiple assignments above and multiple conditions",
			code: []byte(`package main

			func main() {
				c := true
				fooBar := "string"
				a, b := true, false
				if a && b || !c && len(fooBar) > 2 {
					return false
				}

				a, b, c := someFunc()
				if z && len(y) > 0 || 3 == 4 {
					return true
				}
			}`),
			expectedErrorStrings: []string{"if statements should only be cuddled with assigments used in the if statement itself"},
		},
		{
			description: "if statement with multiple assignments, at least one used",
			code: []byte(`package main

			func main() {
				c := true
				fooBar := "string"
				a, b := true, false
				if c && b || !c && len(fooBar) > 2 {
					return false
				}
			}`),
		},
		{
			description: "return can never be cuddleded",
			code: []byte(`package main

			func main() {
				var foo = true
				return
			}`),
			expectedErrorStrings: []string{"return statements should never be cuddled"},
		},
		{
			description: "assigments should only be cuddled with assignments (negative)",
			code: []byte(`package main

			func main() {
				if true {
					fmt.Println("this is bad...")
				}
				foo := true
			}`),
			expectedErrorStrings: []string{"assigments should only be cuddled with other assigments"},
		},
		{
			description: "assigments should only be cuddled with assignments",
			code: []byte(`package main

			func main() {
				bar := false
				foo := true
				biz := true

				return
			}`),
		},
		{
			description: "expressions cannot be cuddled with declarations",
			code: []byte(`package main

			func main() {
				var a bool
				fmt.Println(a)

				return
			}`),
			expectedErrorStrings: []string{"expressions should not be cuddled with declarations or returns"},
		},
		{
			description: "expressions can be cuddlede with assigments",
			code: []byte(`package main

			func main() {
				a := true
				fmt.Println(a)

				return
			}`),
		},
		{
			description: "ranges cannot be cuddled with assigments not used in the range",
			code: []byte(`package main

			func main() {
				anotherList := make([]string, 5)

				myList := make([]int, 10)
				for i := range anotherList {
					fmt.Println(i)
				}
			}`),
			expectedErrorStrings: []string{"ranges should only be cuddled with assignments used in the iteration"},
		},
		{
			description: "ranges can be cuddled if iteration is of assignment above",
			code: []byte(`package main

			func main() {
				myList := make([]string, 5)

				anotherList := make([]int, 10)
				for i := range anotherList {
					fmt.Println(i)
				}
			}`),
		},
		{
			description: "ranges can be cuddled if multiple assignments and/or values iterated",
			code: []byte(`package main

			func main() {
				someList, anotherList := GetTwoListsOfStrings()
				for i := range append(someList, anotherList...) {
					fmt.Println(i)
				}

				aThirdList := GetList()
				for i := range append(someList, aThirdList...) {
					fmt.Println(i)
				}
			}`),
		},
		{
			description: "multi line assignment should not be seen as assignment above",
			code: []byte(`package main

			func main() {
				// This should not be OK since the assignment is split at
				// multiple rows.
				err := SomeFunc(
					"taking multiple values",
					"suddendly not a single line",
				)
				if err != nil {
					fmt.Println("denied")
				}
			}`),
			expectedErrorStrings: []string{"if statements should only be cuddled with assigments"},
		},
		{
			description: "allow cuddle for immediate assignment in block",
			code: []byte(`package main

			func main() {
				// This should be allowed
				idx := i
				if i > 0 {
					idx = i - 1
				}

				vals := map[int]struct{}{}
				for i := range make([]int, 5) {
					vals[i] = struct{}{}
				}

				// This last issue fails.
				x := []int{}

				vals := map[int]struct{}{}
				for i := range make([]int, 5) {
					x = append(x, i)
				}
			}`),
			expectedErrorStrings: []string{"ranges should only be cuddled with assignments used in the iteration"},
		},
		{
			description: "can cuddle only one assignment",
			code: []byte(`package main

			func main() {
				// This is allowed
				foo := true
				bar := false

				biz := true || false
				if biz {
					return true
				}

				// And this is allowed

				foo := true
				bar := false
				biz := true || false

				if biz {
					return true
				}

				// But not this

				foo := true
				bar := false
				biz := true || false
				if biz {
					return false
				}
			}`),
			expectedErrorStrings: []string{"only one cuddle assignment allowed before if statement"},
		},
		{
			description: "using identifies with indicies",
			code: []byte(`package main

			func main() {
				// This should be OK
				runes := []rune{'+', '-'}
				if runes[0] == '+' || runes[0] == '-' {
					return string(runes[1:])
				}

				// And this
				listTwo := []string{}

				listOne := []string{"one"}
				if listOne[0] == "two" {
					return "not allowed"
				}
			}`),
		},
		{
			description: "defer statements can be cuddled",
			code: []byte(`package maino

			import "sync"

			func main() {
				m := sync.Mutex{}

				m.Lock()
				defer m.Unlock()

				// This should (probably?) yield error
				foo := true
				defer func(b bool) {
					fmt.Printf("%v", b)
				}()
			}`),
			expectedErrorStrings: []string{"defer statements should only be cuddled with expressions on same variable"},
		},
		{
			description: "selector expressions are handled like variables",
			code: []byte(`package main

			type T struct {
				List []string
			}

			func main() {
				t := makeT()
				for _, x := range t.List {
					return "this is not like a variable"
				}
			}

			func makeT() T {
				return T{
					List: []string{"one", "two"},
				}
			}`),
			expectedErrorStrings: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r := ProcessFile("unit-test", tc.code)

			require.Len(t, r, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, r[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestTODO(t *testing.T) {
	// All tests added in this TODO section suold fail. To make the
	// implementation easier I'm trying to add the tests and handled them one by
	// one with TDD.
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "append and functions",
			code: []byte(`package main

			func main() {
				var (
					someList = []string{}
				)

				// This should be OK
				foo := true
				someFunc(foo)

				// And this
				bar := "baz"
				someList = append(someList, bar)

				// But not this
				bar := "baz"
				someList = append(someList, "notBar")

				// And not this
				foo := true
				someFunc(false)
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "handle channels (example should succeed)",
			code: []byte(`package main

			func main() {
				timeoutCh := time.After(timeout)

				for range make([]int, 10) {
					select {
					case <-timeoutCh:
						return true
					case <-time.After(10 * time.Millisecond):
						return false
					}
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "handle ForStmt (example should fail)",
			code: []byte(`package main

			func main() {
				bool := true
				for {
					fmt.Println("should not be allowed")

					if bool {
						break
					}
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "trigger stmt *ast.SwitchStmt?!",
			code: []byte(`package main

			func main() {
				t := GetT()

				switch t.GetField() {
				case 1:
					return 0
				case 2:
					return 1
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "support usage if chained",
			code: []byte(`package main

			func main() {
				r := map[string]interface{}{}
				if err := json.NewDecoder(someReader).Decode(&r); err != nil {
					return "this should be OK"
				}
			}`),
			expectedErrorStrings: []string{""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r := ProcessFile("unit-test", tc.code)

			require.Len(t, r, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, r[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

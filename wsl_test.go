package wsl

import (
	"fmt"
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
		{
			description: "do not panic on empty blocks",
			code: []byte(`package main

			func x(r string) (x uintptr)`),
		},
		{
			description: "no false positives for comments",
			code: []byte(`package main

			func main() { // This is a comment, not the first line
				fmt.Println("hello, world")
			}`),
		},
		{
			description: "no false positives for one line functions",
			code: []byte(`package main

			func main() {
				s := func() error { return nil }
			}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			p := NewProcessor()
			p.process("unit-test", tc.code)

			require.Len(t, p.result, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, p.result[i].Reason, tc.expectedErrorStrings[i], "expected error found")
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
			expectedErrorStrings: []string{"block should not start with a whitespace", reasonBlockEndsWithWS},
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
			expectedErrorStrings: []string{"block should not start with a whitespace", reasonBlockEndsWithWS},
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
			description: "Example tests parsed correct when ending with comments",
			code: []byte(`package foo

			func ExampleFoo() {
				fmt.Println("hello world")

				// Output: hello world
			}

			func ExampleBar() {
				fmt.Println("hello world")
				// Output: hello world
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
				reasonBlockEndsWithWS,
				"block should not start with a whitespace",
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
			expectedErrorStrings: []string{reasonNeverCuddleDeclare},
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
			p := NewProcessor()
			p.process("unit-test", tc.code)

			require.Len(t, p.result, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, p.result[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestShouldAddEmptyLines(t *testing.T) {
	cases := []struct {
		description          string
		skip                 bool
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "ok to cuddled if statements directly after assignments multiple values",
			code: []byte(`package main

			func main() {
				v, err := strconv.Atoi("1")
				if err != nil {
					fmt.Println(v)
				}
			}`),
		},
		{
			description: "ok to cuddled if statements directly after assignments single value",
			code: []byte(`package main

			func main() {
				a := true
				if a {
					fmt.Println("I'm OK")
				}
			}`),
		},
		{
			description: "ok to cuddled if statements directly after assignments single value, negations",
			code: []byte(`package main

			func main() {
				a, b := true, false
				if !a {
					fmt.Println("I'm OK")
				}
			}`),
		},
		{
			description: "cannot cuddled if assignments not used",
			code: []byte(`package main

			func main() {
				err := ProduceError()

				a, b := true, false
				if err != nil {
					fmt.Println("I'm OK")
				}
			}`),
			expectedErrorStrings: []string{reasonOnlyCuddleWithUsedAssign},
		},
		{
			description: "cannot cuddled with other things than assignments",
			code: []byte(`package main

			func main() {
				if true {
					fmt.Println("I'm OK")
				}
				if false {
					fmt.Println("I'm OK")
				}
			}`),
			expectedErrorStrings: []string{reasonOnlyCuddleIfWithAssign},
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
			expectedErrorStrings: []string{"only one cuddle assignment allowed before if statement", reasonOnlyCuddleWithUsedAssign},
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

				if foo {
					fmt.Println("return")
				}
				return
			}`),
			expectedErrorStrings: []string{reasonOnlyCuddle2LineReturn},
		},
		{
			description: "assignments should only be cuddled with assignments (negative)",
			code: []byte(`package main

			func main() {
				if true {
					fmt.Println("this is bad...")
				}
				foo := true
			}`),
			expectedErrorStrings: []string{reasonAssignsCuddleAssign},
		},
		{
			description: "assignments should only be cuddled with assignments",
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
			expectedErrorStrings: []string{reasonExpressionCuddledWithDeclOrRet},
		},
		{
			description: "expressions can be cuddlede with assignments",
			code: []byte(`package main

			func main() {
				a := true
				fmt.Println(a)

				return
			}`),
		},
		{
			description: "ranges cannot be cuddled with assignments not used in the range",
			code: []byte(`package main

			func main() {
				anotherList := make([]string, 5)

				myList := make([]int, 10)
				for i := range anotherList {
					fmt.Println(i)
				}
			}`),
			expectedErrorStrings: []string{reasonRangeCuddledWithoutUse},
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
			expectedErrorStrings: []string{reasonRangeCuddledWithoutUse},
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
			description: "simple defer statements",
			code: []byte(`package main

			func main() {
				// OK
				thingOne := getOne()
				thingTwo := getTwo()

				defer thingOne.Close()
				defer thingTwo.Close()

				// OK
				thingOne := getOne()
				defer thingOne.Close()

				thingTwo := getTwo()
				defer thingTwo.Close()

				// NOT OK
				thingOne := getOne()
				defer thingOne.Close()
				thingTwo := getTwo()
				defer thingTwo.Close()

				// NOT OK
				thingOne := getOne()
				thingTwo := getTwo()
				defer thingOne.Close()
				defer thingTwo.Close()
			}`),
			expectedErrorStrings: []string{
				reasonAssignsCuddleAssign,
				reasonOneCuddleBeforeDefer,
				reasonOneCuddleBeforeDefer,
			},
		},
		{
			description: "defer statements can be cuddled",
			code: []byte(`package main

			import "sync"

			func main() {
				m := sync.Mutex{}

				// This should probably be OK
				m.Lock()
				defer m.Unlock()

				foo := true
				defer func(b bool) {
					fmt.Printf("%v", b)
				}()
			}`),
			expectedErrorStrings: []string{
				reasonDeferCuddledWithOtherVar,
			},
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
		},
		{
			description: "return statements may be cuddled if the block is less than three lines",
			code: []byte(`package main

			type T struct {
				Name string
				X	 bool
			}

			func main() {
				t := &T{"someT"}

				t.SetX()
			}

			func (t *T) SetX() *T {
				t.X = true
				return t
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "handle ForStmt",
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
			expectedErrorStrings: []string{reasonForWithoutCondition},
		},
		{
			description: "support usage if chained",
			code: []byte(`package main

			func main() {
				go func() {
					panic("is this real life?")
				}()

				fooFunc := func () {}
				go fooFunc()

				barFunc := func () {}
				go fooFunc()

				go func() {
					fmt.Println("hey")
				}()

				cuddled := true
				go func() {
					fmt.Println("hey")
				}()
			}`),
			expectedErrorStrings: []string{
				reasonGoFuncWithoutAssign,
				reasonGoFuncWithoutAssign,
			},
		},
		{
			description: "switch statements",
			code: []byte(`package main

			func main() {
				// OK
				var b bool
				switch b {
				case true:
					return "a"
				case false:
					return "b"
				}

				// OK
				t := time.Now()

				switch {
				case t.Hour() < 12:
					fmt.Println("It's before noon")
				default:
					fmt.Println("It's after noon")
				}

				// Not ok
				var b bool
				switch anotherBool {
				case true:
					return "a"
				case false:
					return "b"
				}

				// Not ok
				t := time.Now()
				switch {
				case t.Hour() < 12:
					fmt.Println("It's before noon")
				default:
					fmt.Println("It's after noon")
				}
			}`),
			expectedErrorStrings: []string{
				reasonSwitchCuddledWithoutUse,
				reasonAnonSwitchCuddled,
			},
		},
		{
			description: "type switch statements",
			code: []byte(`package main

			func main() {
				// Ok!
				x := GetSome()
				switch v := x.(type) {
				case int:
					return "got int"
				default:
					return "got other"
				}

				// Ok?
				var id string
				switch i := objectID.(type) {
				case int:
					id = strconv.Itoa(i)
				case uint32:
					id = strconv.Itoa(int(i))
				case string:
					id = i
				}

				// Not ok
				var b bool
				switch AnotherVal.(type) {
				case int:
					return "a"
				case string:
					return "b"
				}
			}`),
			expectedErrorStrings: []string{
				reasonTypeSwitchCuddledWithoutUse,
			},
		},
		{
			description: "only cuddle append if appended",
			code: []byte(`package main

			func main() {
				var (
					someList = []string{}
				)

				// Should be OK
				bar := "baz"
				someList = append(someList, bar)

				// But not this
				bar := "baz"
				someList = append(someList, "notBar")
			}`),
			expectedErrorStrings: []string{reasonAppendCuddledWithoutUse},
		},
		{
			description: "cuddle expressions to assignments",
			code: []byte(`package main

			func main() {
				// This should be OK
				foo := true
				someFunc(foo)

				// And not this
				foo := true
				someFunc(false)
			}`),
			expectedErrorStrings: []string{reasonExprCuddlingNonAssignedVar},
		},
		{
			description: "channels and select, no false positives",
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
		},
		{
			description: "switch statements",
			code: []byte(`package main

			func main() {
				t := GetT()
				switch t.GetField() {
				case 1:
					return 0
				case 2:
					return 1
				}

				notT := GetX()
				switch notT {
				case "x":
					return 0
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "select statements",
			code: []byte(`package main

			func main() {
				select {
				case <-time.After(1*time.Second):
					return "1s"
				default:
					return "are we there yet?"
				}
			}`),
		},
		{
			description: "branch statements (continue/break)",
			code: []byte(`package main

			func main() {
				for {
					// Allowed if only one statement in block.
					if true {
						singleLine := true
						break
					}

					// Not allowed for multiple lines
					if true && false {
						multiLine := true
						maybeVar := "var"
						continue
					}

					// Multiple lines require newline
					if false {
						multiLine := true
						maybeVar := "var"

						break
					}
				}
			}`),
			expectedErrorStrings: []string{
				"branch statements should not be cuddled if block has more than two lines",
			},
		},
		{
			description: "append",
			code: []byte(`package main

			func main() {
				var (
					someList = []string{}
				)

				// Should be OK
				bar := "baz"
				someList = append(someList, fmt.Sprintf("%s", bar))

				// Should be OK
				bar := "baz"
				someList = append(someList, []string{"foo", bar}...)

				// Should be OK
				bar := "baz"
				someList = append(someList, bar)

				// But not this
				bar := "baz"
				someList = append(someList, "notBar")

				// Ok
				biz := "str"
				whatever := ItsJustAppend(biz)

				// Ok
				zzz := "str"
				whatever := SoThisIsOK(biz)
			}`),
			expectedErrorStrings: []string{
				reasonAppendCuddledWithoutUse,
			},
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
		},
		{
			description: "support function literals",
			code: []byte(`package main

			func main() {
				n := 1
				f := func() {

					// This violates whitespaces

					// Violates cuddle
					notThis := false
					if isThisTested {
						return

					}
				}

				f()
			}`),
			expectedErrorStrings: []string{
				"block should not start with a whitespace",
				reasonBlockEndsWithWS,
				reasonOnlyCuddleWithUsedAssign,
			},
		},
		{
			description: "locks",
			code: []byte(`package main

			func main() {
				hashFileCache.Lock()
				out, ok := hashFileCache.m[file]
				hashFileCache.Unlock()

				mu := &sync.Mutex{}
				mu.X(y).Z.RLock()
				x, y := someMap[someKey]
				mu.RUnlock()
			}`),
		},
		{
			description: "append type",
			code: []byte(`package main

			func main() {
				s := []string{}

				// Multiple append should be OK
				s = append(s, "one")
				s = append(s, "two")
				s = append(s, "three")

				// Both assigned and called should be allowed to be cuddled.
				// It's not nice, but they belong'
				p.defs = append(p.defs, x)
				def.parseFrom(p)
				p.defs = append(p.defs, def)
				def.parseFrom(p)
				def.parseFrom(p)
				p.defs = append(p.defs, x)
			}`),
		},
		{
			description: "AllowAssignAndCallCuddle",
			code: []byte(`package main

			func main() {
				for i := range make([]int, 10) {
					fmt.Println("x")
				}
				x.Calling()
			}`),
			expectedErrorStrings: []string{
				reasonExpressionCuddledWithBlock,
			},
		},
		{
			description: "indexes in maps and arrays",
			code: []byte(`package main

			func main() {
				for i := range make([]int, 10) {
					key := GetKey()
					if val, ok := someMap[key]; ok {
						fmt.Println("ok!")
					}

					someOtherMap := GetMap()
					if val, ok := someOtherMap[key]; ok {
						fmt.Println("ok")
					}

					someIndex := 3
					if val := someSlice[someIndex]; val != nil {
						retunr
					}
				}
			}`),
		},
		{
			description: "splice slice, concat key",
			code: []byte(`package main

			func main() {
				start := 0
				if v := aSlice[start:3]; v {
					fmt.Println("")
				}

				someKey := "str"
				if v, ok := someMap[obtain(someKey)+"str"]; ok {
					fmt.Println("Hey there")
				}

				end := 10
				if v := arr[3:notEnd]; !v {
					// Error
				}

				notKey := "str"
				if v, ok := someMap[someKey]; ok {
					// Error
				}
			}`),
			expectedErrorStrings: []string{
				reasonOnlyCuddleWithUsedAssign,
				reasonOnlyCuddleWithUsedAssign,
			},
		},
		{
			description: "multiline case statements",
			code: []byte(`package main

			func main() {
				switch {
				case true,
					false:
					fmt.Println("ok")
				case true ||
					false:
					fmt.Println("ok")
				case true, false:
					fmt.Println("ok")
				case true || false:
					fmt.Println("ok")
				case true,
					false:

					fmt.Println("starting whitespace multiline case")
				case true ||
					false:

					fmt.Println("starting whitespace multiline case")
				case true,
					false:
					fmt.Println("ending whitespace multiline case")

				case true,
					false:
					fmt.Println("last case should also fail")

				}
			}`),
			expectedErrorStrings: []string{
				reasonBlockEndsWithWS,
				reasonBlockStartsWithWS,
				reasonBlockStartsWithWS,
			},
		},
		{
			description: "allow http body close best practice",
			code: []byte(`package main

			func main() {
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
			}`),
		},
		{
			description: "multiple go statements can be cuddled",
			code: []byte(`package main

			func main() {
				t1 := NewT()
				t2 := NewT()
				t3 := NewT()

				go t1()
				go t2()
				go t3()

				multiCuddle1 := NewT()
				multiCuddle2 := NewT()
				go multiCuddle2()

				multiCuddle3 := NeT()
				go multiCuddle1()

				t4 := NewT()
				t5 := NewT()
				go t5()
				go t4()
			}`),
			expectedErrorStrings: []string{
				reasonOneCuddleBeforeGo,
				reasonGoFuncWithoutAssign,
				reasonOneCuddleBeforeGo,
			},
		},
		{
			description: "boilerplate",
			code: []byte(`package main

			func main() {
				counter := 0
				if somethingTrue {
					counter++
				}

				counterTwo := 0
				if somethingTrue {
					counterTwo--
				}

				notCounter := 0
				if somethingTrue {
					counter--
				}
			}`),
			expectedErrorStrings: []string{
				reasonOnlyCuddleWithUsedAssign,
			},
		},
		{
			description: "ensure blocks are found and checked for errors",
			code: []byte(`package main

				func main() {
					var foo string

					x := func() {
						var err error
						foo = "1"
					}()

					x := func() {
						var err error
						foo = "1"
					}

					func() {
						var err error
						foo = "1"
					}()

					func() {
						var err error
						foo = "1"
					}

					func() {
						func() {
							return func() {
								var err error
								foo = "1"
							}
						}()
					}

					var x error
					foo, err := func() { return "", nil }

					defer func() {
						var err error
						foo = "1"
					}()

					go func() {
						var err error
						foo = "1"
					}()
				}`),
			expectedErrorStrings: []string{
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
				reasonAssignsCuddleAssign,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			if tc.skip {
				t.Skip("not implemented")
			}

			p := NewProcessor()
			p.process("unit-test", tc.code)

			require.Len(t, p.result, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, p.result[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestWithConfig(t *testing.T) {
	// This test is used to simulate code to perform TDD. This part should never
	// be committed with any test.
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
		customConfig         *Configuration
	}{
		{
			description: "AllowAssignAndCallCuddle",
			code: []byte(`package main

			func main() {
				p.token(':')
				d.StartBit = p.uint()
				p.token('|')
				d.Size = p.uint()
			}`),
			customConfig: &Configuration{
				AllowAssignAndCallCuddle: true,
			},
		},
		{
			description: "multi line assignment ok when enabled",
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
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle: false,
			},
			expectedErrorStrings: []string{reasonOnlyCuddleIfWithAssign},
		},
		{
			description: "multi line assignment not ok when disabled",
			code: []byte(`package main

			func main() {
				// This should not be OK since the assignment is split at
				// multiple rows.
				err := SomeFunc(
					"taking multiple values",
					"suddenly not a single line",
				)
				if err != nil {
					fmt.Println("accepted")
				}

				// This should still fail since it's not the same variable
				noErr := SomeFunc(
					"taking multiple values",
					"suddenly not a single line",
				)
				if err != nil {
					fmt.Println("rejected")
				}
			}`),
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle: true,
			},
			expectedErrorStrings: []string{reasonOnlyCuddleWithUsedAssign},
		},
		{
			description: "strict append",
			code: []byte(`package main

			func main() {
				x := []string{}
				y := "not in x"
				x = append(x, "not y") // Should not be allowed

				z := "in x"
				x = append(x, z) // Should be allowed

				m := transform(z)
				x := append(x, z) // Should be allowed
			}`),
			customConfig: &Configuration{
				StrictAppend: true,
			},
			expectedErrorStrings: []string{reasonAppendCuddledWithoutUse},
		},
		{
			description: "allow cuddle var",
			code: []byte(`package main

			func main() {
				var t bool
				var err error

				var t = true
				if t {
					fmt.Println("x")
				}
			}`),
			customConfig: &Configuration{
				AllowCuddleDeclaration: true,
			},
		},
		{
			description: "support to end blocks with a comment",
			code: []byte(`package main

			func main() {
				for i := range make([]int, 10) {
					fmt.Println("x")
					// This is OK
				}

				for i := range make([]int, 10) {
					fmt.Println("x")
					// Pad
					// With empty lines
					//
					// This is OK
				}

				for i := range make([]int, 10) {
					fmt.Println("x")

					// This is NOT OK
				}

				for i := range make([]int, 10) {
					fmt.Println("x")
					// This is NOT OK

				}

				for i := range make([]int, 10) {
					fmt.Println("x")
					/* Block comment one line OK */
				}

				for i := range make([]int, 10) {
					fmt.Println("x")
					/*
						Block comment one multiple lines OK
					*/
				}

				switch {
				case true:
					fmt.Println("x")
					// This is OK
				case false:
					fmt.Println("y")
					// This is OK
				}
			}`),
			customConfig: &Configuration{
				AllowTrailingComment: true,
			},
			expectedErrorStrings: []string{
				reasonBlockEndsWithWS,
				reasonBlockEndsWithWS,
			},
		},
		{
			description: "case blocks can end with or without newlines and comments",
			code: []byte(`package main

			func main() {
				switch "x" {
				case "a correct": // Test
					fmt.Println("1")
					fmt.Println("1")
				case "a correct":
					fmt.Println("1")
					fmt.Println("2")

				case "b correct":
					fmt.Println("1")
				case "b correct":
					fmt.Println("1")

				case "b bad":
					fmt.Println("1")
					// I'm not allowed'
				case "b bad":
					fmt.Println("1")
					// I'm not allowed'

				case "no false positives for if":
					fmt.Println("checking")

					// Some comment
					if 1 < 2 {
						// Another comment
						if 2 > 1 {
							fmt.Println("really sure")
						}
					}
				case "end": // This is a case
					fmt.Println("4")
				}
			}`),
			customConfig: &Configuration{
				ForceCaseTrailingWhitespaceLimit: 0,
				AllowTrailingComment:             false,
			},
		},
		{
			description: "enforce newline at size 3",
			code: []byte(`package main

			func main() {
				switch "x" {
				case "a correct": // Test
					fmt.Println("1")
				case "a ok":
					fmt.Println("1")

				case "b correct":
					fmt.Println("1")
					fmt.Println("1")
					fmt.Println("1")

				case "b bad":
					fmt.Println("1")
					fmt.Println("1")
					fmt.Println("1")
				case "b correct":
					fmt.Println("1")
					fmt.Println("1")
					fmt.Println("1")
					// This is OK

				case "b bad":
					fmt.Println("1")
					fmt.Println("1")
					fmt.Println("1")
					// This is not OK, no whitespace
				case "b ok":
					fmt.Println("1")
					fmt.Println("1")
					fmt.Println("1")
					/*
						This is not OK
						Multiline but no whitespace
					*/
				case "end": // This is a case
					fmt.Println("4")
				}
			}`),
			customConfig: &Configuration{
				ForceCaseTrailingWhitespaceLimit: 3,
				AllowTrailingComment:             true,
			},
			expectedErrorStrings: []string{
				reasonCaseBlockTooCuddly,
				reasonCaseBlockTooCuddly,
				reasonCaseBlockTooCuddly,
			},
		},
		{
			description: "must cuddle error checks with the error assignment",
			customConfig: &Configuration{
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				err := errors.New("bad thing")

				if err != nil {
					fmt.Println("I'm OK")
				}
			}`),
			expectedErrorStrings: []string{reasonMustCuddleErrCheck},
		},
		{
			description: "must cuddle error checks with the error assignment multivalue",
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle:   true,
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				b, err := FSByte(useLocal, name)

				if err != nil {
					panic(err)
				}
			}`),
			expectedErrorStrings: []string{reasonMustCuddleErrCheck},
		},
		{
			description: "must cuddle error checks with the error assignment only on assignment",
			customConfig: &Configuration{
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				var (
					err error
					once sync.Once
				)

				once.Do(func() {
					err = ProduceError()
				})
			
				if err != nil {
					return nil, err
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "must cuddle error checks with the error assignment multivalue NoError",
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle:   true,
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				b, err := FSByte(useLocal, name)
				if err != nil {
					panic(err)
				}
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "must cuddle error checks with the error assignment known err",
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle:   true,
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				result, err := FetchSomething()

				if err == sql.ErrNoRows {
					return []Model{}
				}
			}`),
			expectedErrorStrings: []string{reasonMustCuddleErrCheck},
		},
		{
			description: "err-check cuddle enforcement doesn't generate false-positives.",
			customConfig: &Configuration{
				AllowMultiLineAssignCuddle:   true,
				ForceCuddleErrCheckAndAssign: true,
				ErrorVariableNames:           []string{"err"},
			},
			code: []byte(`package main

			import "errors"

			func main() {
				result, err := FetchSomething()
				if err == sql.ErrNoRows {
					return []Model{}
				}

				foo := generateFoo()
				if foo == "bar" {
 				    handleBar()
                }

				var baz []string

				err = loadStuff(&baz)
				if err != nil{
					return nil
				}

				if err := ProduceError(); err != nil {
					return err
				}

				return baz
			}`),
			expectedErrorStrings: []string{},
		},
		{
			description: "allow separated leading comment",
			customConfig: &Configuration{
				AllowSeparatedLeadingComment: true,
			},
			code: []byte(`package main

			func main() {
				// These blocks should not generate error
				func () {
					// Comment

					// Comment
					fmt.Println("Hello, World")
				}

				func () {
					/*
						Multiline
					*/

					/*
						Multiline
					*/
					fmt.Println("Hello, World")
				}

				func () {
					/*
						Multiline
					*/

					// Comment
					fmt.Println("Hello, World")
				}

				func () {
					// Comment

					/*
						Multiline
					*/
					fmt.Println("Hello, World")
				}

				func () { // Comment
					/*
						Multiline
					*/
					fmt.Println("Hello, World")
				}
			}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			p := NewProcessor()
			if tc.customConfig != nil {
				p = NewProcessorWithConfig(*tc.customConfig)
			}

			p.process("unit-test", tc.code)
			require.Len(t, p.result, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, p.result[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

func TestTODO(t *testing.T) {
	// This test is used to simulate code to perform TDD. This part should never
	// be committed with any test.
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
		customConfig         *Configuration
	}{
		{
			description: "boilerplate",
			code: []byte(`package main

			func main() {
				for i := range make([]int, 10) {
					fmt.Println("x")
				}
			}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			p := NewProcessor()
			if tc.customConfig != nil {
				p = NewProcessorWithConfig(*tc.customConfig)
			}

			p.process("unit-test", tc.code)

			t.Logf("WARNINGS: %s", p.warnings)

			for _, r := range p.result {
				fmt.Println(r.String())
			}

			require.Len(t, p.result, len(tc.expectedErrorStrings), "correct amount of errors found")

			for i := range tc.expectedErrorStrings {
				assert.Contains(t, p.result[i].Reason, tc.expectedErrorStrings[i], "expected error found")
			}
		})
	}
}

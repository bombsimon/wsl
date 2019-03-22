package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidFile(t *testing.T) {
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

func TestIgnoreComments(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "comments are ignored, ignore this crappy code",
			code: []byte(`package main

			/*
			Multi line comments ignore errors
			if true {

				fmt.Println("Well hello blankies!")

			}
			*/

			// Also signel comments are ignored
			// foo := false
			// if true {
			//
			//	foo = true
			//
			// }`),
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

func TestNoEmptyStart(t *testing.T) {
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

func TestBlankLineAfterBracket(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
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
			expectedErrorStrings: []string{"declarations can never be cuddled"},
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
				assert.Contains(t, r[i].Reason, tc.expectedErrorStrings[i], fmt.Sprintf("expected error found at index %d", i))
			}
		})
	}
}

func TestBlankLineBeforeIf(t *testing.T) {
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
			expectedErrorStrings: []string{"if statements can only be cuddled with assigments used in the if statement itself"},
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
			expectedErrorStrings: []string{"if statements can only be cuddled with assigments"},
		},
	}

	/*
		TODO: Checks that allows this:
		foo := true
		bar := false

		biz := true || false
		if biz {
		}

		AND

		foo := true
		bar := false
		biz := true || false

		if biz {
		}

		BUT NOT THIS

		foo := true
		bar := false
		biz := true || false
		if biz {
		}
	*/

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

func TestCuddle(t *testing.T) {
	cases := []struct {
		description          string
		code                 []byte
		expectedErrorStrings []string
	}{
		{
			description: "return can never be cuddleded",
			code: []byte(`package main

			func main() {
				var foo = true
				return
			}`),
			expectedErrorStrings: []string{"return statements can never be cuddled"},
		},
		{
			description: "assigments can only be cuddled with assignments (negative)",
			code: []byte(`package main

			func main() {
				if true {
					fmt.Println("this is bad...")
				}
				foo := true
			}`),
			expectedErrorStrings: []string{"assigments can only be cuddled with other assigments"},
		},
		{
			description: "assigments can only be cuddled with assignments",
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
			expectedErrorStrings: []string{"expressions can not be cuddled with decarations or returns"},
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

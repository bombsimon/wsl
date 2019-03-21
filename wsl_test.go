package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIgnoreComments(t *testing.T) {
	code := []byte(`package main

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
// }`)

	r := ProcessFile("unit-test", code)
	assert.Len(t, r, 0, "no errors for comments")
}

func TestNoEmptyStart(t *testing.T) {
	code := []byte(`func a() {
	fmt.Println("this is function a")
}`)

	r := ProcessFile("unit-test", code)

	require.Len(t, r, 1, "invalid syntax, cannot be linted")
	assert.Equal(t, r[0].Reason, "invalid syntax, file cannot be linted")
}

/*
	var (
		l      []string
		r      []Result
		assert = assert.New(t)
	)

	l = []string{
		"func a() {",
		"",
		"	var a = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(2, r[0].LineNo, "Error found at second line")

	l = []string{
		"var notFirstLine = 0",
		"",
		"if true {",
		"",
		"	// true is true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(4, r[0].LineNo, "Error found at fourth line")
}

func TestBlankLineAfterBracket(t *testing.T) {
	var (
		l      []string
		r      []Result
		assert = assert.New(t)
	)

	l = []string{
		"var b = true",
		"",
		"if b {",
		"	// b is true here",
		"}",
		"var a = false",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(6, r[0].LineNo, "Error found at fifth line")

	l = []string{
		"if true {",
		"	if false {",
		"		if true && false {",
		"			// Whelp",
		"		}",
		"	}",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No errors when nestling if statements")

	l = []string{
		"if true {",
		"	if false {",
		"		if true && false {",
		"			// Whelp",
		"		}",
		"		thisIsToTight = true",
		"	}",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(6, r[0].LineNo, "Error found at sixth line")
}

func TestBlankLineBeforeIf(t *testing.T) {
	var (
		l      []string
		r      []Result
		assert = assert.New(t)
	)

	l = []string{
		"_, err := strconv.Atoi(\"1\")",
		"if err != nil {",
		"	// ok to cuddle err checks",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No errors when cuddling the most common err check")

	l = []string{
		"var a = true",
		"if a {",
		"	// I'm cuddling!",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No error when using one assign value")

	l = []string{
		"foo, bar, baz := SomeFunc()",
		"if !baz {",
		"	var a = true",
		"}",
	}

	l = []string{
		"result := DoIt()",
		"if err != nil {",
		"	// Error?!",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(2, r[0].LineNo, "Error found at second line")

	l = []string{
		"var a = true",
		"if a {",
		"	// I'm cuddling!",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No error when using one of any assign value")

	l = []string{
		"foo := true",
		"bar := false",
		"",
		"biz := Some()",
		"if !biz {",
		"	var a = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No errors using room before last assignment")

	l = []string{
		"foo := true",
		"bar := false",
		"biz := Some()",
		"",
		"if !biz {",
		"	var a = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(0, len(r), "No errors using room before if statement")

	l = []string{
		"foo := true",
		"bar := false",
		"biz := Some()",
		"if !biz {",
		"	var a = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")

	l = []string{
		"switch {",
		"case 1:",
		"",
		"	// I start with a blank line",
		"	notAllowed = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
}

func TestFirstLineBlank(t *testing.T) {
	var (
		l      []string
		r      []Result
		assert = assert.New(t)
	)

	l = []string{
		"",
		"package main",
		"",
		"func main() {",
		"	fmt.Println(\"Hello world\")",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(1, len(r), "An error was found")
	assert.Equal(1, r[0].LineNo, "Error found at first line")
}

func TestFileProcessing(t *testing.T) {
	var (
		r      []Result
		assert = assert.New(t)
	)

	r = ProcessFiles([]string{"testfiles/01.go"})
	assert.Equal(8, len(r), "An error was found")
	assert.Equal(9, r[0].LineNo, "Error was found on line 4")

	r = ProcessFile("testfiles/02.go")
	assert.Equal(0, len(r), "No errors in file")
}
*/

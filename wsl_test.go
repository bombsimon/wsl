package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEmptyStart(t *testing.T) {
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
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 2, "Error found at second line")

	l = []string{
		"var notFirstLine = 0",
		"",
		"if true {",
		"",
		"	// true is true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 4, "Error found at fourth line")
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
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 6, "Error found at fifth line")

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
	assert.Equal(len(r), 0, "No errors when nestling if statements")

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
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 6, "Error found at sixth line")
}

func TestBlankLineBeforeIf(t *testing.T) {
	var (
		l      []string
		r      []Result
		assert = assert.New(t)
	)

	l = []string{
		"var a = true",
		"if a {",
		"	// I'm cuddling!",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 2, "Error found at second line")

	l = []string{
		"_, err := strconv.Atoi(\"1\")",
		"if err != nil {",
		"	// ok to cuddle err checks",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(len(r), 0, "No errors when cuddling err check")

	l = []string{
		"switch {",
		"case 1:",
		"",
		"	// I start with a blank line",
		"	notAllowed = true",
		"}",
	}

	r = ProcessLines(l, "nofile")
	assert.Equal(len(r), 1, "An error was found")
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
	assert.Equal(len(r), 1, "An error was found")
	assert.Equal(r[0].LineNo, 1, "Error found at first line")
}

func TestFileProcessing(t *testing.T) {
	var (
		r      []Result
		assert = assert.New(t)
	)

	r = ProcessFile("testfiles/01.go")
	assert.Equal(len(r), 8, "An error was found")
	assert.Equal(r[0].LineNo, 9, "Error was found on line 4")

	r = ProcessFile("testfiles/02.go")
	assert.Equal(len(r), 0, "No errors in file")
}

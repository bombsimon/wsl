package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	whitespaceIndicator = "}"
)

// Result represnets the result of one linted row which was not accepted by the
// white space linter.
type Result struct {
	FileName string
	LineNo   int
	Reason   string
}

// ProcessFiles takes a string slice with file names (full paths) and lints
// them.
func ProcessFiles(files []string) []Result {
	var result []Result

	for _, file := range files {
		result = append(result, ProcessFile(file)...)
	}

	return result
}

// ProcessFile will read a single file and lint it.
func ProcessFile(file string) []Result {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	return ProcessLines(lines, file)
}

// ProcessLines will process a string slice line by line to determine if
// there's enough room in the code to make it readable.
func ProcessLines(lines []string, filename string) []Result {
	var (
		result    []Result
		inComment bool
	)

	for i := range lines {
		idx := i

		if i > 0 {
			idx = i - 1
		}

		var (
			lineNo      = i + 1
			currentRow  = strings.TrimSpace(lines[i])
			previousRow = strings.TrimSpace(lines[idx])
		)

		// Check if we're at the end of a block comment.
		if strings.HasPrefix(currentRow, "*/") {
			inComment = false
			continue
		}

		// Check if we're at the start of a block comment.
		if strings.HasPrefix(currentRow, "/*") {
			inComment = true
		}

		// Allow whatever crappy style in comments.
		if inComment {
			continue
		}

		if strings.HasPrefix(currentRow, "//") {
			continue
		}

		if i == 0 {
			if currentRow == "" {
				result = append(result, Result{filename, lineNo, "first line should never be blank"})
			}

			continue
		}

		if previousRow == whitespaceIndicator {
			if currentRow != "" && currentRow != ")" && currentRow != whitespaceIndicator && !strings.HasPrefix(currentRow, "case") {
				// Closing bracket (}) followed by zero or more parentheses and maybe a comma.
				// Should see '})' and '},' as valid lines without newline after '}'.
				matched, _ := regexp.MatchString(`}\)*,?$`, currentRow)

				if !matched {
					result = append(result, Result{filename, lineNo, fmt.Sprintf("must be a blank line after '%s'", whitespaceIndicator)})
				}
			}
		}

		if previousRow == "" {
			if currentRow == whitespaceIndicator {
				result = append(result, Result{filename, lineNo, fmt.Sprintf("blank line not allowed before '%s'", whitespaceIndicator)})
			}
		}

		// Don't allow the start of a block to be an empty line.
		if strings.HasSuffix(previousRow, "{") || strings.HasPrefix(previousRow, "case") {
			if currentRow == "" {
				result = append(result, Result{filename, lineNo, "blank line not allowed in beginning of a block"})
			}
		}

		if strings.HasPrefix(currentRow, "if") && !emptyOrComment(previousRow) {
			// We can (should) cuddle if previous line was if or case
			if !strings.HasSuffix(previousRow, "{") && !strings.HasPrefix(previousRow, "case") {
				assigned := getAssigned(previousRow)

				if !inList(currentRow, assigned) {
					result = append(result, Result{filename, lineNo, "if statement should have a blank line before they start, unless they use variables assigned on the line before"})
				} else {
					if lineNo >= 3 {
						twoUp := strings.TrimSpace(lines[i-2])

						if strings.Contains(twoUp, "=") {
							result = append(result, Result{filename, lineNo - 2, "if statements can check variables assigned on the line above only if the assignment row follows a line without assignment"})
						}
					}
				}
			}
		}
	}

	return result
}

func inList(s string, list []string) bool {
	currentRowWords := strings.Split(s, " ")

	// Iterate all words on this row (if !foo || bar > 0)
	for _, w := range currentRowWords {
		// Compare towards list with assigned on previous line (baz, bar)
		for _, l := range list {
			// if !foo  == baz || !foo == !baz
			// if bar == bar || bar == !bar
			if w == l || w == "!"+l {
				return true
			}
		}
	}

	return false
}

func getAssigned(s string) []string {
	var assigned []string

	parts := strings.Split(s, " ")

	for _, part := range parts {
		// If the token is an assignment, no more variables to store.
		if part == "=" || part == ":=" {
			break
		}

		assigned = append(assigned, part)
	}

	return assigned
}

func emptyOrComment(s string) bool {
	return s == "" || strings.HasPrefix(s, "//")
}

package wsl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

// ProcessDirectory will read a directory (recursive if desired) and process
// each file found in passed directory/directories.
func ProcessDirectory(dir string, recursive bool) []Result {
	// Always ignore .git directory.
	_, f := filepath.Split(dir)
	if f == ".git" {
		return []Result{}
	}

	var result []Result

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileOrDir := filepath.Join(dir, file.Name())

		if s, err := os.Stat(fileOrDir); err == nil && s.IsDir() {
			if recursive {
				result = append(result, ProcessDirectory(fileOrDir, recursive)...)
			}

			continue
		}

		result = append(result, ProcessFile(fileOrDir)...)
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
	var result []Result

	for i := range lines {
		idx := i

		if i > 0 {
			idx = i - 1
		}

		var (
			lineNo             = i + 1
			currentRow         = lines[i]
			currentRowTrimmed  = strings.TrimSpace(currentRow)
			previousRow        = lines[idx]
			previousRowTrimmed = strings.TrimSpace(previousRow)
		)

		if i == 0 {
			if currentRowTrimmed == "" {
				result = append(result, Result{filename, lineNo, "first line should never be blank"})
			}

			continue
		}

		// Don't lint comments.
		if strings.HasPrefix(currentRowTrimmed, "//") {
			continue
		}

		if previousRowTrimmed == whitespaceIndicator {
			if currentRow != "" && currentRowTrimmed != whitespaceIndicator {
				result = append(result, Result{filename, lineNo, fmt.Sprintf("must be a blank line after '%s'", whitespaceIndicator)})
			}
		}

		if previousRowTrimmed == "" {
			if currentRowTrimmed == whitespaceIndicator {
				result = append(result, Result{filename, lineNo, fmt.Sprintf("blank line not allowed before '%s'", whitespaceIndicator)})
			}
		}

		// Don't allow the start of a block to be an empty line.
		if strings.HasSuffix(previousRowTrimmed, "{") {
			if currentRowTrimmed == "" {
				result = append(result, Result{filename, lineNo, "blank line not allowed in beginning of a block"})
			}
		}

		if strings.HasPrefix(currentRowTrimmed, "if") {
			if !hasIfPrefix(currentRowTrimmed, "err", "ok", "!ok") && !strings.HasSuffix(previousRowTrimmed, "{") && !emptyOrComment(previousRowTrimmed) {
				result = append(result, Result{filename, lineNo, "if statement should have a blank line before they start, unless for error checking or nested"})
			}
		}
	}

	return result
}

func emptyOrComment(s string) bool {
	return s == "" || strings.HasPrefix(s, "//")
}

func hasIfPrefix(s string, ss ...string) bool {
	for _, prefix := range ss {
		if strings.HasPrefix(s, fmt.Sprintf("if %s", prefix)) {
			return true
		}
	}

	return false
}

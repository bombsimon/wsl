package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bombsimon/wsl"
)

func main() {
	var (
		args         []string
		help         bool
		noTest       bool
		showWarnings bool
		cwd, _       = os.Getwd()
		files        = []string{}
		finalFiles   = []string{}
		config       = wsl.DefaultConfig()
	)

	flag.BoolVar(&help, "h", false, "Show this help text")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&noTest, "n", false, "Don't lint test files")
	flag.BoolVar(&noTest, "no-test", false, "")
	flag.BoolVar(&showWarnings, "w", false, "Show warnings (ignored rules)")
	flag.BoolVar(&showWarnings, "warnings", false, "")

	flag.BoolVar(&config.CheckAppend, "check-append", true, "Strict rules for append")
	flag.BoolVar(&config.AllowAssignAndCallsCuddle, "allow-assign-and-call", true, "Allow assignments and calls to be cuddled (if using same variable/type)")
	flag.BoolVar(&config.AllowMultiLineAssignmentCuddled, "allow-multi-line-assignment", true, "Allow cuddling with multi line assignments")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	args = flag.Args()
	if len(args) == 0 {
		args = []string{"./..."}
	}

	for _, f := range args {
		if strings.HasSuffix(f, "/...") {
			dir, _ := filepath.Split(f)

			files = append(files, expandGoWildcard(dir)...)

			continue
		}

		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}

	// Use relative path to print shorter names, sort out test files if chosen.
	for _, f := range files {
		if noTest {
			if strings.HasSuffix(f, "_test.go") {
				continue
			}
		}

		if relativePath, err := filepath.Rel(cwd, f); err == nil {
			finalFiles = append(finalFiles, relativePath)

			continue
		}

		finalFiles = append(finalFiles, f)
	}

	processor := wsl.NewProcessorWithConfig(config)
	result, warnings := processor.ProcessFiles(finalFiles)

	for _, r := range result {
		fmt.Println(r.String())
	}

	if showWarnings && len(warnings) > 0 {
		fmt.Println()
		fmt.Println("⚠️  Warnings found")

		for _, w := range warnings {
			fmt.Println(w)
		}
	}

	if len(result) > 0 {
		os.Exit(2)
	}
}

func expandGoWildcard(root string) []string {
	foundFiles := []string{}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Only append go files
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		foundFiles = append(foundFiles, path)

		return nil
	})

	return foundFiles
}

func showHelp() {
	helpText := `Usage: wsl <file> [files...]

Also supports package syntax but will use it in relative path, i.e. ./pkg/...

Flags:`

	fmt.Println(helpText)
	flag.PrintDefaults()
}

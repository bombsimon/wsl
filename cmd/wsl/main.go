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
		args       []string
		help       bool
		notest     bool
		cwd, _     = os.Getwd()
		files      = []string{}
		finalFiles = []string{}
	)

	flag.BoolVar(&help, "h", false, "Show this help text")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&notest, "n", false, "Don't lint test files")
	flag.BoolVar(&notest, "no-test", false, "")

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
		if notest {
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

	r := wsl.ProcessFiles(finalFiles)

	for _, x := range r {
		fmt.Printf("%s:%d: %s\n", x.FileName, x.LineNumber, x.Reason)
	}

	if len(r) > 0 {
		os.Exit(2)
	}

	return
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

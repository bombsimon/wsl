package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bombsimon/ff"
)

func main() {
	var (
		recursive bool
	)

	flag.BoolVar(&recursive, "recursive", false, "parse file recursive")
	flag.BoolVar(&recursive, "r", false, "parse file recursive")

	flag.Parse()

	var (
		files []string
		err   error
		match []ff.Match
	)

	switch flag.NArg() {
	case 0:
		match, err = ff.FilesFromPattern(".", "*.go", "", recursive)
	case 1:
		match, err = ff.FilesFromPattern(".", flag.Args()[0], "", recursive)
	default:
		match, err = ff.FilesFromPattern(flag.Args()[0], flag.Args()[1], "", recursive)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, m := range match {
		files = append(files, m.FullName())
	}

	r := ProcessFiles(files)

	for _, x := range r {
		fmt.Printf("%s:%d: %s\n", x.FileName, x.LineNo, x.Reason)
	}

	if len(r) > 0 {
		os.Exit(2)
	}

	os.Exit(0)
}

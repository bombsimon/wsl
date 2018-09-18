package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bombsimon/wsl"
)

func main() {
	recursive := flag.Bool("recursive", false, "parse file recursive")

	flag.Parse()

	var r []wsl.Result

	switch flag.NArg() {
	case 0:
		r = append(r, wsl.ProcessDirectory(".", *recursive)...)
	default:
		for _, files := range flag.Args() {
			files, _ := filepath.Glob(files)

			for _, name := range files {
				if s, err := os.Stat(name); err == nil && s.IsDir() {
					r = append(r, wsl.ProcessDirectory(name, *recursive)...)
				} else {
					r = append(r, wsl.ProcessFile(name)...)
				}
			}
		}
	}

	for _, x := range r {
		fmt.Printf("%s:%d: %s\n", x.FileName, x.LineNo, x.Reason)
	}

	if len(r) > 0 {
		os.Exit(1)
	}

	os.Exit(0)
}

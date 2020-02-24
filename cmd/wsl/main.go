package main

import (
	"github.com/bombsimon/wsl/v3"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wsl.Analyzer)
}

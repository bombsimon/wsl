package main

import (
	"github.com/bombsimon/wsl/v2"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wsl.Analyzer)
}

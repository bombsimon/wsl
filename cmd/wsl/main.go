package main

import (
	"github.com/bombsimon/wsl/v5"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wsl.NewAnalyzer(nil))
}

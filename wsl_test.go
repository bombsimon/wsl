package wsl

import (
	"flag"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFixer(t *testing.T) {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	analyzer := Analyzer
	analyzer.Flags = *flags

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "wip")
}

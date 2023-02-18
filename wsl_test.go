package wsl

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFixer(t *testing.T) {
	analyzer := Analyzer
	analyzer.Flags = flags()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "fixer")
}

func TestDefaultConfig(t *testing.T) {
	analyzer := Analyzer
	analyzer.Flags = flags()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "default_config")
}

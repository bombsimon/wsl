package wsl

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(nil)

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "default_config")
}

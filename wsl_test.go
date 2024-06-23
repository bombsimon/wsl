package wsl

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(&Configuration{})

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "default_config")
}

func TestWithConfig(t *testing.T) {
	testdata := analysistest.TestData()

	for _, tc := range []struct {
		subdir   string
		configFn func(*Configuration)
	}{
		{
			subdir: "if_errcheck",
			configFn: func(config *Configuration) {
				config.Errcheck = true
			},
		},
	} {
		t.Run(tc.subdir, func(t *testing.T) {
			config := &Configuration{}
			tc.configFn(config)

			analyzer := NewAnalyzer(config)
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, filepath.Join("with_config", tc.subdir))
		})
	}
}

func TestWIP(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(&Configuration{})

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "wip")
}

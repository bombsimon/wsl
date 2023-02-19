package wsl

import (
	"flag"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWIP(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = flags()

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "wip")
}

func TestDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = flags()

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "default_config")
}

func TestWithConfig(t *testing.T) {
	defaultConfig := config

	flags := flag.NewFlagSet("", flag.ExitOnError)

	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = *flags

	for _, tc := range []struct {
		subdir   string
		configFn func()
	}{
		{
			subdir: "case_blocks",
			configFn: func() {
				config.ForceCaseTrailingWhitespaceLimit = 3
			},
		},
		{
			subdir: "multi_line_assign",
			configFn: func() {
				config.AllowMultiLineAssignCuddle = false
			},
		},
		{
			subdir: "assign_and_call",
			configFn: func() {
				config.AllowAssignAndCallCuddle = false
			},
		},
		{
			subdir: "trailing_comments",
			configFn: func() {
				config.AllowTrailingComment = false
			},
		},
	} {
		t.Run(tc.subdir, func(t *testing.T) {
			config = defaultConfig
			tc.configFn()

			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, filepath.Join("with_config", tc.subdir))
		})
	}
}

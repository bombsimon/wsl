package wsl

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWIP(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(nil)

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "wip")
}

func TestDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(nil)

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "default_config")
}

func TestWithConfig(t *testing.T) {
	testdata := analysistest.TestData()

	for _, tc := range []struct {
		subdir   string
		configFn func(*Configuration)
	}{
		{
			subdir: "case_blocks",
			configFn: func(config *Configuration) {
				config.ForceCaseTrailingWhitespaceLimit = 3
			},
		},
		{
			subdir: "multi_line_assign",
			configFn: func(config *Configuration) {
				config.AllowMultiLineAssignCuddle = false
			},
		},
		{
			subdir: "assign_and_call",
			configFn: func(config *Configuration) {
				config.AllowAssignAndCallCuddle = false
			},
		},
		{
			subdir: "trailing_comments",
			configFn: func(config *Configuration) {
				config.AllowTrailingComment = true
			},
		},
		{
			subdir: "separate_leading_whitespace",
			configFn: func(config *Configuration) {
				config.AllowSeparatedLeadingComment = true
			},
		},
		{
			subdir: "error_check",
			configFn: func(config *Configuration) {
				config.ForceCuddleErrCheckAndAssign = true
			},
		},
		{
			subdir: "short_decl",
			configFn: func(config *Configuration) {
				config.ForceExclusiveShortDeclarations = true
			},
		},
		{
			subdir: "strict_append",
			configFn: func(config *Configuration) {
				config.StrictAppend = false
			},
		},
		{
			subdir: "assign_and_anything",
			configFn: func(config *Configuration) {
				config.AllowAssignAndAnythingCuddle = true
			},
		},
		{
			subdir: "decl",
			configFn: func(config *Configuration) {
				config.AllowCuddleDeclaration = true
			},
		},
		{
			subdir: "include_generated",
			configFn: func(config *Configuration) {
				config.IncludeGenerated = true
			},
		},
	} {
		t.Run(tc.subdir, func(t *testing.T) {
			config := defaultConfig()
			tc.configFn(config)

			analyzer := NewAnalyzer(config)
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, filepath.Join("with_config", tc.subdir))
		})
	}
}

package wsl

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultConfig(t *testing.T) {
	dirs, err := os.ReadDir("./testdata/src/default_config")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range dirs {
		t.Run(tc.Name(), func(t *testing.T) {
			testdata := analysistest.TestData()
			analyzer := NewAnalyzer(NewConfig())

			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, filepath.Join("default_config", tc.Name()))
		})
	}
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
				config.Checks.Add(CheckErr)
			},
		},
		{
			subdir: "no_check_decl",
			configFn: func(config *Configuration) {
				config.Checks.Remove(CheckDecl)
			},
		},
		{
			subdir: "whole_block",
			configFn: func(config *Configuration) {
				config.AllowWholeBlock = true
			},
		},
		{
			subdir: "first_in_block_n1",
			configFn: func(config *Configuration) {
				config.AllowFirstInBlock = true
				config.FirstInBlockMaxDepth = 1
			},
		},
		{
			subdir: "first_in_block_n5",
			configFn: func(config *Configuration) {
				config.AllowFirstInBlock = true
				config.FirstInBlockMaxDepth = 5
			},
		},
	} {
		t.Run(tc.subdir, func(t *testing.T) {
			config := NewConfig()
			tc.configFn(config)

			analyzer := NewAnalyzer(config)
			analysistest.RunWithSuggestedFixes(t, testdata, analyzer, filepath.Join("with_config", tc.subdir))
		})
	}
}

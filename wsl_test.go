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
				config.Checks[CheckErr] = struct{}{}
			},
		},
		{
			subdir: "no_check_decl",
			configFn: func(config *Configuration) {
				delete(config.Checks, CheckDecl)
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

package wsl

import (
	"flag"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFixer(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = flags()

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "fixer")
}

func TestDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = flags()

	analysistest.Run(t, testdata, analyzer, "default_config")
}

func TestWithConfig(t *testing.T) {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	flags.BoolVar(&config.StrictAppend, "strict-append", true, "Strict rules for append")
	flags.BoolVar(&config.AllowAssignAndCallCuddle, "allow-assign-and-call", false, "Allow assignments and calls to be cuddled (if using same variable/type)")
	flags.BoolVar(&config.AllowAssignAndAnythingCuddle, "allow-assign-and-anything", false, "Allow assignments and anything to be cuddled")
	flags.BoolVar(&config.AllowMultiLineAssignCuddle, "allow-multi-line-assign", false, "Allow cuddling with multi line assignments")
	flags.BoolVar(&config.AllowCuddleDeclaration, "allow-cuddle-declarations", false, "Allow declarations to be cuddled")
	flags.BoolVar(&config.AllowTrailingComment, "allow-trailing-comment", false, "Allow blocks to end with a comment")
	flags.BoolVar(&config.AllowSeparatedLeadingComment, "allow-separated-leading-comment", false, "Allow empty newlines in leading comments")
	flags.BoolVar(&config.ForceCuddleErrCheckAndAssign, "force-err-cuddling", false, "Force cuddling of error checks with error var assignment")
	flags.BoolVar(&config.ForceExclusiveShortDeclarations, "force-short-decl-cuddling", false, "Force short declarations to cuddle by themselves")
	flags.IntVar(&config.ForceCaseTrailingWhitespaceLimit, "force-case-trailing-whitespace", 0, "Force newlines for case blocks > this number.")

	testdata := analysistest.TestData()
	analyzer := Analyzer
	analyzer.Flags = *flags

	analysistest.Run(t, testdata, analyzer, "with_config")
}

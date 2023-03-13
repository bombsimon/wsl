package wsl

import (
	"flag"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the wsl analyzer.
//
//nolint:gochecknoglobals // Analyzers tend to be global.
var Analyzer = &analysis.Analyzer{
	Name:             "wsl",
	Doc:              "add or remove empty lines",
	Flags:            flags(),
	Run:              run,
	RunDespiteErrors: true,
}

//nolint:gochecknoglobals // Config needs to be overwritten with flags.
var config = Configuration{
	AllowAssignAndAnythingCuddle:     false,
	AllowAssignAndCallCuddle:         true,
	AllowCuddleDeclaration:           false,
	AllowMultiLineAssignCuddle:       true,
	AllowSeparatedLeadingComment:     false,
	AllowTrailingComment:             false,
	ForceCuddleErrCheckAndAssign:     false,
	ForceExclusiveShortDeclarations:  false,
	StrictAppend:                     true,
	AllowCuddleWithCalls:             []string{"Lock", "RLock"},
	AllowCuddleWithRHS:               []string{"Unlock", "RUnlock"},
	ErrorVariableNames:               []string{"err"},
	ForceCaseTrailingWhitespaceLimit: 0,
}

func flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	flags.BoolVar(&config.AllowAssignAndAnythingCuddle, "allow-assign-and-anything", false, "Allow assignments and anything to be cuddled")
	flags.BoolVar(&config.AllowAssignAndCallCuddle, "allow-assign-and-call", true, "Allow assignments and calls to be cuddled (if using same variable/type)")
	flags.BoolVar(&config.AllowCuddleDeclaration, "allow-cuddle-declarations", false, "Allow declarations to be cuddled")
	flags.BoolVar(&config.AllowMultiLineAssignCuddle, "allow-multi-line-assign", true, "Allow cuddling with multi line assignments")
	flags.BoolVar(&config.AllowSeparatedLeadingComment, "allow-separated-leading-comment", false, "Allow empty newlines in leading comments")
	flags.BoolVar(&config.AllowTrailingComment, "allow-trailing-comment", false, "Allow blocks to end with a comment")
	flags.BoolVar(&config.ForceCuddleErrCheckAndAssign, "force-err-cuddling", false, "Force cuddling of error checks with error var assignment")
	flags.BoolVar(&config.ForceExclusiveShortDeclarations, "force-short-decl-cuddling", false, "Force short declarations to cuddle by themselves")
	flags.BoolVar(&config.StrictAppend, "strict-append", true, "Strict rules for append")
	flags.IntVar(&config.ForceCaseTrailingWhitespaceLimit, "force-case-trailing-whitespace", 0, "Force newlines for case blocks > this number.")

	flags.String("allow-cuddle-with-calls", "Lock,RLock", "Comma separated list of idents that can have cuddles after")
	flags.String("allow-cuddle-with-rhs", "Unlock,RUnlock", "Comma separated list of idents that can have cuddles before")
	flags.String("error-variable-names", "err", "Comma separated list of error variable names")

	return *flags
}

func lookupAndSplit(pass *analysis.Pass, flagName string) []string {
	flagVal := pass.Analyzer.Flags.Lookup(flagName)
	if flagVal == nil {
		return nil
	}

	return strings.Split(flagVal.Value.String(), ",")
}

func run(pass *analysis.Pass) (interface{}, error) {
	if v := lookupAndSplit(pass, "allow-cuddle-with-calls"); v != nil {
		config.AllowCuddleWithCalls = v
	}

	if v := lookupAndSplit(pass, "allow-cuddle-with-rhs"); v != nil {
		config.AllowCuddleWithRHS = v
	}

	if v := lookupAndSplit(pass, "error-variable-names"); v != nil {
		config.ErrorVariableNames = v
	}

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		processor := newProcessorWithConfig(file, pass.Fset, &config)
		processor.parseAST()

		for pos, fix := range processor.result {
			textEdits := []analysis.TextEdit{}
			for _, f := range fix.fixRanges {
				textEdits = append(textEdits, analysis.TextEdit{
					Pos:     f.fixRangeStart,
					End:     f.fixRangeEnd,
					NewText: f.newText,
				})
			}

			pass.Report(analysis.Diagnostic{
				Pos:      pos,
				Category: "whitespace",
				Message:  fix.reason,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						TextEdits: textEdits,
					},
				},
			})
		}
	}

	//nolint:nilnil // A pass don't need to return anything.
	return nil, nil
}

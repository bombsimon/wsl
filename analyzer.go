package wsl

import (
	"flag"

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
	StrictAppend:                     true,
	AllowAssignAndCallCuddle:         true,
	AllowMultiLineAssignCuddle:       true,
	AllowTrailingComment:             false,
	ForceCuddleErrCheckAndAssign:     true,
	ForceCaseTrailingWhitespaceLimit: 0,
	AllowCuddleWithCalls:             []string{"Lock", "RLock"},
	AllowCuddleWithRHS:               []string{"Unlock", "RUnlock"},
	ErrorVariableNames:               []string{"err"},
}

func flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	flags.BoolVar(&config.StrictAppend, "strict-append", true, "Strict rules for append")
	flags.BoolVar(&config.AllowAssignAndCallCuddle, "allow-assign-and-call", true, "Allow assignments and calls to be cuddled (if using same variable/type)")
	flags.BoolVar(&config.AllowAssignAndAnythingCuddle, "allow-assign-and-anything", false, "Allow assignments and anything to be cuddled")
	flags.BoolVar(&config.AllowMultiLineAssignCuddle, "allow-multi-line-assign", true, "Allow cuddling with multi line assignments")
	flags.BoolVar(&config.AllowCuddleDeclaration, "allow-cuddle-declarations", false, "Allow declarations to be cuddled")
	flags.BoolVar(&config.AllowTrailingComment, "allow-trailing-comment", false, "Allow blocks to end with a comment")
	flags.BoolVar(&config.AllowSeparatedLeadingComment, "allow-separated-leading-comment", false, "Allow empty newlines in leading comments")
	flags.BoolVar(&config.ForceCuddleErrCheckAndAssign, "force-err-cuddling", false, "Force cuddling of error checks with error var assignment")
	flags.BoolVar(&config.ForceExclusiveShortDeclarations, "force-short-decl-cuddling", false, "Force short declarations to cuddle by themselves")
	flags.IntVar(&config.ForceCaseTrailingWhitespaceLimit, "force-case-trailing-whitespace", 0, "Force newlines for case blocks > this number.")

	return *flags
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		processor := NewProcessorWithConfig(file, pass.Fset, &config)
		processor.ParseAST()

		for _, v := range processor.Result {
			var (
				newText   []byte
				textEdits []analysis.TextEdit
			)

			switch v.Type {
			case WhitespaceShouldAddBefore:
				newText = []byte("\n")
			case WhitespaceShouldAddAfter:
				newText = []byte("\n")
			case WhitespaceShouldRemoveEnd:
				newText = []byte("\n}")
			case WhitespaceShouldRemoveBeginning:
				newText = []byte("{\n")
			default:
				//nolint:gocritic // Not while TODO
				// TODO
				continue
			}

			textEdits = []analysis.TextEdit{
				{
					Pos:     v.FixRangeStart,
					End:     v.FixRangeEnd,
					NewText: newText,
				},
			}

			d := analysis.Diagnostic{
				Pos:      v.ReportAt,
				Category: "",
				Message:  v.Reason,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message:   v.Type.String(),
						TextEdits: textEdits,
					},
				},
			}

			pass.Report(d)
		}
	}

	//nolint:nilnil // A pass don't need to return anything.
	return nil, nil
}

package wsl

import (
	"flag"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the wsl analyzer.
// nolint // gochecknoglobals - Not needed
var Analyzer = &analysis.Analyzer{
	Name:             "wsl",
	Doc:              "add or remove empty lines",
	Flags:            flags(),
	Run:              run,
	RunDespiteErrors: true,
}

var config = Configuration{
	StrictAppend:                     true,
	AllowAssignAndCallCuddle:         true,
	AllowMultiLineAssignCuddle:       true,
	AllowTrailingComment:             false,
	ForceCuddleErrCheckAndAssign:     false,
	ForceCaseTrailingWhitespaceLimit: 0,
	AllowCuddleWithCalls:             []string{"Lock", "RLock"},
	AllowCuddleWithRHS:               []string{"Unlock", "RUnlock"},
	ErrorVariableNames:               []string{"err"},
}

func flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	flags.IntVar(&config.ForceCaseTrailingWhitespaceLimit, "force-case-trailing-whitespace", 0, "Force newlines for case blocks > this number.")

	return *flags
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		processor := NewProcessorWithConfig(file, pass.Fset, config)
		processor.ParseAST()

		for _, v := range processor.Result {
			var (
				pos       token.Pos
				end       token.Pos
				newText   []byte
				textEdits []analysis.TextEdit
			)

			switch v.Type {
			case WhitespaceShouldAddBefore:
				pos = v.Node.Pos()
				end = v.Node.Pos()
				newText = []byte("\n")
			case WhitespaceShouldAddAfter:
				pos = v.Node.End()
				end = v.Node.End()
				newText = []byte("\n")
			default:
				// TODO
				continue
			}

			if !v.NoFix {
				textEdits = []analysis.TextEdit{
					{
						Pos:     pos,
						End:     end,
						NewText: newText,
					},
				}
			}

			d := analysis.Diagnostic{
				Pos:      v.Node.Pos(),
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

	return nil, nil
}

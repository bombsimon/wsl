package wsl

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the wsl analyzer.
// nolint // gochecknoglobals - Not needed
var Analyzer = &analysis.Analyzer{
	Name:             "wsl",
	Doc:              "add or remove empty lines",
	Run:              run,
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		processor := NewProcessor(file, pass.Fset)
		processor.ParseAST()

		for _, v := range processor.Result {
			var (
				pos     token.Pos
				end     token.Pos
				newText []byte
			)

			switch v.Type {
			case WhitespaceShouldAdd:
				pos = v.Node.Pos()
				end = v.Node.Pos()
				newText = []byte("\n")
			default:
				// TODO
				continue
			}

			d := analysis.Diagnostic{
				Pos:      v.Node.Pos(),
				Category: "",
				Message:  v.Reason,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: v.Type.String(),
						TextEdits: []analysis.TextEdit{
							{
								Pos:     pos,
								End:     end,
								NewText: newText,
							},
						},
					},
				},
			}

			pass.Report(d)
		}
	}

	return nil, nil
}

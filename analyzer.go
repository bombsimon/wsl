package wsl

import (
	"flag"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer(config *Configuration) *analysis.Analyzer {
	wa := &wslAnalyzer{config: config}

	return &analysis.Analyzer{
		Name:             "wsl",
		Doc:              "add or remove empty lines",
		Flags:            wa.flags(),
		Run:              wa.run,
		RunDespiteErrors: true,
	}
}

// wslAnalyzer is a wrapper around the configuration which is used to be able to
// set the configuration when creating the analyzer and later be able to update
// flags and running method.
type wslAnalyzer struct {
	config *Configuration
}

func (wa *wslAnalyzer) flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	return *flags
}

func (wa *wslAnalyzer) run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		wsl := New(file, pass, wa.config)
		wsl.Run()

		for pos, fix := range wsl.Issues {
			textEdits := []analysis.TextEdit{}

			for _, f := range fix.FixRanges {
				textEdits = append(textEdits, analysis.TextEdit{
					Pos:     f.FixRangeStart,
					End:     f.FixRangeEnd,
					NewText: []byte("\n"),
				})
			}

			pass.Report(analysis.Diagnostic{
				Pos:      pos,
				Category: "whitespace",
				Message:  fix.Message,
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

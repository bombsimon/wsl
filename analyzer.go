package wsl

import (
	"flag"
	"go/ast"
	"go/token"
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
//
// TODO: This should be a public configuration you can pass to `NewAnalyzer` to
// enable `golang-lint` to pass `includeGenerated` but also pass a simpler
// enable/disable list of strings.
type wslAnalyzer struct {
	config           *Configuration
	includeGenerated bool
	enableAll        bool
	disableAll       bool
	enable           []string
	disable          []string
}

func (wa *wslAnalyzer) flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	// If we have a configuration set we're not running from the command line so
	// we don't use any flags.
	if wa.config != nil {
		return *flags
	}

	wa.config = NewConfig()

	flags.BoolVar(&wa.includeGenerated, "include-generated", false, "Include generated files")
	flags.BoolVar(&wa.enableAll, "enable-all", false, "Enable all checks")
	flags.BoolVar(&wa.disableAll, "disable-all", false, "Disable all checks")
	flags.Var(&multiStringValue{slicePtr: &wa.enable}, "enable", "Comma separated list of checks to enable")
	flags.Var(&multiStringValue{slicePtr: &wa.disable}, "disable", "Comma separated list of checks to disable")

	return *flags
}

func (wa *wslAnalyzer) run(pass *analysis.Pass) (any, error) {
	if err := wa.config.Update(wa.enableAll, wa.disableAll, wa.enable, wa.disable); err != nil {
		return nil, err
	}

	for _, file := range pass.Files {
		filename := getFilename(pass.Fset, file)
		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		// if the file is related to cgo the filename of the unadjusted position
		// is a not a '.go' file.
		unadjustedFilename := pass.Fset.PositionFor(file.Pos(), false).Filename

		// if the file is related to cgo the filename of the unadjusted position
		// is a not a '.go' file.
		if !strings.HasSuffix(unadjustedFilename, ".go") {
			continue
		}

		// The file is skipped if the "unadjusted" file is a Go file, and it's a
		// generated file (ex: "_test.go" file). The other non-Go files are
		// skipped by the first 'if' with the adjusted position.
		if !wa.includeGenerated && ast.IsGenerated(file) {
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

// multiStringValue is a flag that supports multiple values. It's implemented to
// contain a pointer to a string slice that will be overwritten when the flag's
// `Set` method is called.
type multiStringValue struct {
	slicePtr *[]string
}

// Set implements the flag.Value interface and will overwrite the pointer to
// the
// slice with a new pointer after splitting the flag by comma.
func (m *multiStringValue) Set(value string) error {
	var s []string

	for v := range strings.SplitSeq(value, ",") {
		s = append(s, strings.TrimSpace(v))
	}

	*m.slicePtr = s

	return nil
}

// Set implements the flag.Value interface.
func (m *multiStringValue) String() string {
	if m.slicePtr == nil {
		return ""
	}

	return strings.Join(*m.slicePtr, ", ")
}

func getFilename(fset *token.FileSet, file *ast.File) string {
	filename := fset.PositionFor(file.Pos(), true).Filename
	if !strings.HasSuffix(filename, ".go") {
		return fset.PositionFor(file.Pos(), false).Filename
	}

	return filename
}

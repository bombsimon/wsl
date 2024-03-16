package wsl

import (
	"go/ast"
	"go/token"
)

const (
	MessageAddWhitespace    = "missing whitespace decreases readability"
	MessageRemoveWhitespace = "unnecessary whitespace decreases readability"
)

type Configuration struct{}

type FixRange struct {
	FixRangeStart token.Pos
	FixRangeEnd   token.Pos
}

type Issue struct {
	Message   string
	FixRanges []FixRange // TODO: Do we need multiple?
}

type WSL struct {
	File     *ast.File
	Fset     *token.FileSet
	Comments ast.CommentMap
	Issues   map[token.Pos]Issue // TODO: When do we have multiple reports?
}

func New(file *ast.File, fset *token.FileSet, _ *Configuration) *WSL {
	return &WSL{
		Fset:   fset,
		File:   file,
		Issues: make(map[token.Pos]Issue),
	}
}

func (w *WSL) Run() {
	for _, decl := range w.File.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			w.CheckFuncDecl(funcDecl)
		}
	}
}

func (w *WSL) CheckFuncDecl(funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	w.CheckUnnecessaryBlockLeadingNewline(funcDecl.Body)
}

func (w *WSL) CheckUnnecessaryBlockLeadingNewline(body *ast.BlockStmt) {
	// No statements in the block, let's leave it as is.
	if len(body.List) == 0 {
		return
	}

	openingPos := body.Lbrace + 1
	firstStmt := body.List[0].Pos()

	comments := ast.NewCommentMap(w.Fset, body, w.File.Comments)
	for _, commentGroup := range comments {
		for _, comment := range commentGroup {
			// The comment starts after the current opening position (originally
			// the LBrace) and ends before the current first statement
			// (originally first body.List item).
			if comment.Pos() > openingPos && comment.End() < firstStmt {
				openingPosLine := w.lineFor(openingPos)
				commentStartLine := w.lineFor(comment.Pos())

				// If comment starts at the same line as the opening position it
				// should just extend the position for the fixer if needed.
				// func fn() { // This comment starts at the same line as LBrace
				switch {
				// The comment is on the same line as current opening position.
				// E.g. func fn() { // A comment
				case commentStartLine == openingPosLine:
					openingPos = comment.End()
				// Opening position is the same as `{` and the comment is
				// directly on the line after (no empty line)
				case openingPosLine == w.lineFor(body.Lbrace) &&
					commentStartLine == w.lineFor(body.Lbrace)+1:
					openingPos = comment.End()
				// The opening position has been updated, it's another comment.
				case openingPosLine != w.lineFor(body.Lbrace):
					openingPos = comment.End()
				// The opening position is still { and the comment is not
				// directly above - it must be an empty line which shouldn't be
				// there.
				default:
					firstStmt = comment.Pos()
				}
			}
		}
	}

	openingPosLine := w.Fset.PositionFor(openingPos, false).Line
	firstStmtLine := w.Fset.PositionFor(firstStmt, false).Line

	if openingPosLine != firstStmtLine-1 {
		w.addError(openingPos, openingPos, firstStmt, MessageRemoveWhitespace)
	}
}

func (w *WSL) lineFor(pos token.Pos) int {
	return w.Fset.PositionFor(pos, false).Line
}

func (w *WSL) addError(report, start, end token.Pos, message string) {
	issue, ok := w.Issues[report]
	if !ok {
		issue = Issue{
			Message:   message,
			FixRanges: []FixRange{},
		}
	}

	issue.FixRanges = append(issue.FixRanges, FixRange{
		FixRangeStart: start,
		FixRangeEnd:   end,
	})

	w.Issues[report] = issue
}

package wsl

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/analysis"
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
	TypeInfo *types.Info
	Comments ast.CommentMap
	Issues   map[token.Pos]Issue // TODO: When do we have multiple reports?
}

func New(file *ast.File, pass *analysis.Pass, _ *Configuration) *WSL {
	return &WSL{
		Fset:     pass.Fset,
		File:     file,
		TypeInfo: pass.TypesInfo,
		Issues:   make(map[token.Pos]Issue),
	}
}

func (w *WSL) Run() {
	for _, decl := range w.File.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			w.CheckFunc(funcDecl)
		}
	}
}

func (w *WSL) CheckFunc(funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	w.CheckUnnecessaryBlockLeadingNewline(funcDecl.Body)
	w.CheckBlock(funcDecl.Body)
}

func (w *WSL) CheckIf(stmt *ast.IfStmt, cursor *Cursor) {
	current := cursor.currentIdx

	n := w.numberOfStatementsAbove(cursor)
	if n > 0 {
		// TODO: Check if we're cuddled and if so that we do it as we should. Rules:
		// * Only cuddle with _one_ assignment
		// * Only cuddled with an assignment declared on the line above
		// * OR a variable that is used first in the block
		// * OR - if configured - anywhere in the block.
	}

	// TODO: If we're _not_ cuddled, check if the previous variable implements
	// an error and if that's used in the if statement.
	currentIdents := allIdents(cursor.Stmt())

	shouldCuddleErr := true // TODO(config): Should be configurable
	if shouldCuddleErr && n == 0 && cursor.Previous() {
		previousIdents := allIdents(cursor.Stmt())
		intersects := identIntersection(currentIdents, previousIdents)

		if slices.ContainsFunc(intersects, func(ident *ast.Ident) bool {
			return w.implementsErr(ident)
		}) {
			w.addError(
				stmt.Pos(),
				cursor.Stmt().End(),
				stmt.Pos(),
				MessageRemoveWhitespace,
			)

			// If we add the error at the same position but with a different fix
			// range, only the fix range will be updated.
			//
			//   a := 1
			//   err := fn()
			//
			//   if err != nil {}
			//
			// Should become
			//
			//   a := 1
			//
			//   err := fn()
			//   if err != nil {}
			if w.numberOfStatementsAbove(cursor) > 0 {
				w.addError(
					stmt.Pos(),
					cursor.Stmt().Pos(),
					cursor.Stmt().Pos(),
					MessageRemoveWhitespace,
				)
			}
		}
	}

	cursor.currentIdx = current

	// if
	w.CheckBlock(stmt.Body)

	switch v := stmt.Else.(type) {
	// else-if
	case *ast.IfStmt:
		w.CheckIf(v, cursor)
	// else
	case *ast.BlockStmt:
		w.CheckBlock(v)
	}
}

func (w *WSL) CheckBlock(block *ast.BlockStmt) {
	cursor := NewCursor(0, block.List)
	for cursor.Next() {
		fmt.Printf("%d: %T\n", cursor.currentIdx, cursor.Stmt())

		//nolint:gocritic // This is not commented out code, it's examples
		switch s := cursor.Stmt().(type) {
		// if a {} else if b {} else {}
		case *ast.IfStmt:
			w.CheckIf(s, cursor)
		// for {} / for a; b; c {}
		case *ast.ForStmt:
		// for _, _ = range a {}
		case *ast.RangeStmt:
		// switch {} // switch a {}
		case *ast.SwitchStmt:
		// switch a.(type) {}
		case *ast.TypeSwitchStmt:
		// return a
		case *ast.ReturnStmt:
		// continue / break
		case *ast.BranchStmt:
		// var a
		case *ast.DeclStmt:
		// a := a
		case *ast.AssignStmt:
		// a++ / a--
		case *ast.IncDecStmt:
		// defer func() {}
		case *ast.DeferStmt:
		// go func() {}
		case *ast.GoStmt:
		case *ast.ExprStmt:
		default:
			fmt.Printf("Not implemented: %T\n", s)
		}
	}
}

// numberOfStatementsAbove will find out how many lines above the cursor's
// current statement there is without any newlines between.
func (w *WSL) numberOfStatementsAbove(cursor *Cursor) int {
	cursor.Save()
	defer cursor.Reset()

	statementsWithoutNewlines := 0
	currentStmtStartLine := w.lineFor(cursor.Stmt().Pos())

	for cursor.Previous() {
		previousStmtEndLine := w.lineFor(cursor.Stmt().End())
		if previousStmtEndLine != currentStmtStartLine-1 {
			break
		}

		currentStmtStartLine = w.lineFor(cursor.Stmt().Pos())
		statementsWithoutNewlines++
	}

	return statementsWithoutNewlines
}

func (w *WSL) CheckUnnecessaryBlockLeadingNewline(body *ast.BlockStmt) {
	// No statements in the block, let's leave it as is.
	if len(body.List) == 0 {
		return
	}

	lbraceLine := w.lineFor(body.Lbrace)
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
				case openingPosLine == lbraceLine &&
					commentStartLine == lbraceLine+1:
					openingPos = comment.End()
				// The opening position has been updated, it's another comment.
				case openingPosLine != lbraceLine:
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

func (w *WSL) implementsErr(node *ast.Ident) bool {
	typeInfo := w.TypeInfo.TypeOf(node)

	errorType, ok := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	if !ok {
		return false
	}

	return types.Implements(typeInfo, errorType)
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

func allIdents(node ast.Node) []*ast.Ident {
	idents := []*ast.Ident{}

	if node == nil {
		return idents
	}

	switch n := node.(type) {
	case *ast.Ident:
		return []*ast.Ident{n}
	case *ast.AssignStmt:
		for _, lhs := range n.Lhs {
			idents = append(idents, allIdents(lhs)...)
		}
	case *ast.IfStmt:
		idents = append(idents, allIdents(n.Cond)...)
		// TODO: idents = append(idents, allIdents(n.Else)...)
	case *ast.BinaryExpr:
		idents = append(idents, allIdents(n.X)...)
		idents = append(idents, allIdents(n.Y)...)
	case *ast.BasicLit:
	default:
		spew.Dump(node)
		fmt.Printf("%T\n", node)
	}

	return idents
}

func identIntersection(a, b []*ast.Ident) []*ast.Ident {
	intersects := []*ast.Ident{}

	for _, as := range a {
		for _, bs := range b {
			if as.Name == bs.Name {
				intersects = append(intersects, as)
			}
		}
	}

	return intersects
}

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

// CheckType is a type that represents a checker to run.
type CheckType int

// Each checker is represented by a CheckType that is used to enable or disable
// the check.
const (
	CheckAssign CheckType = iota
	CheckBreak
	CheckContinue
	CheckDecl
	CheckDefer
	CheckExpr
	CheckFor
	CheckGo
	CheckIf
	CheckLeadingWhitespace
	CheckRange
	CheckReturn
	CheckSwitch
	CheckTypeSwitch
)

type Configuration struct {
	// Require no newline between error assignment and error check.
	Errcheck bool
	Checks   map[CheckType]struct{}
}

func NewConfig() *Configuration {
	return &Configuration{
		Errcheck: false,
		Checks: map[CheckType]struct{}{
			CheckAssign:            {},
			CheckBreak:             {},
			CheckContinue:          {},
			CheckDecl:              {},
			CheckDefer:             {},
			CheckExpr:              {},
			CheckFor:               {},
			CheckGo:                {},
			CheckIf:                {},
			CheckLeadingWhitespace: {},
			CheckRange:             {},
			CheckReturn:            {},
			CheckSwitch:            {},
			CheckTypeSwitch:        {},
		},
	}
}

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
	Config   *Configuration
}

func New(file *ast.File, pass *analysis.Pass, cfg *Configuration) *WSL {
	return &WSL{
		Fset:     pass.Fset,
		File:     file,
		TypeInfo: pass.TypesInfo,
		Issues:   make(map[token.Pos]Issue),
		Config:   cfg,
	}
}

func (w *WSL) Run() {
	for _, decl := range w.File.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			w.CheckFunc(funcDecl)
		}
	}
}

func (w *WSL) CheckCuddling(stmt ast.Node, cursor *Cursor, maxAllowedStatements int) {
	reset := cursor.Save()
	defer reset()

	currentIdents := allIdents(cursor.Stmt())
	previousIdents := []*ast.Ident{}

	var previousNode ast.Node

	if cursor.Previous() {
		previousNode = cursor.Stmt()
		previousIdents = allIdents(previousNode)

		cursor.Next() // Move forward again
	}

	_, prevIsAssign := previousNode.(*ast.AssignStmt)
	_, prevIsDecl := previousNode.(*ast.DeclStmt)
	_, currIsDefer := stmt.(*ast.DeferStmt)
	n := w.numberOfStatementsAbove(cursor)

	if n > 0 {
		if prevIsAssign || prevIsDecl || currIsDefer {
			intersects := identIntersection(currentIdents, previousIdents)

			// No idents above share name with one in the if statement.
			if len(intersects) == 0 {
				w.addError(
					stmt.Pos(),
					stmt.Pos(),
					stmt.Pos(),
					MessageAddWhitespace,
				)
			}

			// Check the maximum number of allowed statements above, and if not
			// disabled (-1) check that the previous one intersects with the
			// current one.
			if maxAllowedStatements != -1 && n > maxAllowedStatements {
				// Idents on the line above exist in the current condition so that
				// should remain cuddled.
				if len(intersects) > 0 {
					w.addError(
						previousNode.Pos(),
						previousNode.Pos(),
						previousNode.Pos(),
						MessageAddWhitespace,
					)
				}
				// TODO: Features:
				// * Allow idents that is used first in the block
				// * OR - if configured - anywhere in the block.
			}
		} else {
			// If we have a statement above and it's not an assignment or
			// declaration we unconditionally add an error.
			w.addError(
				cursor.Stmt().Pos(),
				cursor.Stmt().Pos(),
				cursor.Stmt().Pos(),
				MessageAddWhitespace,
			)
		}
	}

	if _, ok := stmt.(*ast.IfStmt); !ok {
		return
	}

	if w.Config.Errcheck && n == 0 && len(previousIdents) > 0 {
		if slices.ContainsFunc(previousIdents, func(ident *ast.Ident) bool {
			return w.implementsErr(ident)
		}) {
			w.addError(
				stmt.Pos(),
				previousNode.End(),
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
			cursor.Previous()

			if w.numberOfStatementsAbove(cursor) > 0 {
				w.addError(
					stmt.Pos(),
					previousNode.Pos(),
					previousNode.Pos(),
					MessageAddWhitespace,
				)
			}
		}
	}
}

func (w *WSL) CheckFunc(funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	w.CheckBlock(funcDecl.Body)
}

func (w *WSL) CheckIf(stmt *ast.IfStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckIf]; ok {
		w.CheckCuddling(stmt, cursor, 1)
	}

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

func (w *WSL) CheckFor(stmt *ast.ForStmt, cursor *Cursor) {
	defer w.CheckBlock(stmt.Body)

	if _, ok := w.Config.Checks[CheckFor]; !ok {
		return
	}

	w.CheckCuddling(stmt, cursor, 1)
}

func (w *WSL) CheckRange(stmt *ast.RangeStmt, cursor *Cursor) {
	defer w.CheckBlock(stmt.Body)

	if _, ok := w.Config.Checks[CheckRange]; !ok {
		return
	}

	w.CheckCuddling(stmt, cursor, 1)
}

func (w *WSL) CheckSwitch(stmt *ast.SwitchStmt, cursor *Cursor) {
	defer w.CheckBlock(stmt.Body)

	if _, ok := w.Config.Checks[CheckSwitch]; !ok {
		return
	}

	w.CheckCuddling(stmt, cursor, 1)
}

func (w *WSL) CheckTypeSwitch(stmt *ast.TypeSwitchStmt, cursor *Cursor) {
	defer w.CheckBlock(stmt.Body)

	if _, ok := w.Config.Checks[CheckTypeSwitch]; !ok {
		return
	}

	w.CheckCuddling(stmt, cursor, 1)
}

func (w *WSL) CheckExprStmt(stmt *ast.ExprStmt, cursor *Cursor) {
	defer w.CheckExpr(stmt.X, cursor)

	if _, ok := w.Config.Checks[CheckExpr]; !ok {
		return
	}

	previousNode := cursor.PreviousNode()

	// We can cuddle any amount call statements so only check cuddling if the
	// previous one isn't a function call.
	if _, ok := previousNode.(*ast.ExprStmt); !ok {
		w.CheckCuddling(stmt, cursor, -1)
	}
}

func (w *WSL) CheckGo(stmt *ast.GoStmt, cursor *Cursor) {
	defer w.CheckExpr(stmt.Call, cursor)

	if _, ok := w.Config.Checks[CheckGo]; !ok {
		return
	}

	previousNode := cursor.PreviousNode()

	// We can cuddle any amount `go` statements so only check cuddling if the
	// previous one isn't a `go` call.
	// We don't even have to check if it's actually cuddled or just the previous
	// one because even if it's not but is a `go` statement it's valid.
	if _, ok := previousNode.(*ast.GoStmt); !ok {
		w.CheckCuddling(stmt, cursor, 1)
	}
}

func (w *WSL) CheckDefer(stmt *ast.DeferStmt, cursor *Cursor) {
	defer w.CheckExpr(stmt.Call, cursor)

	if _, ok := w.Config.Checks[CheckDefer]; !ok {
		return
	}

	previousNode := cursor.PreviousNode()

	// We can cuddle any amount `defer` statements so only check cuddling if the
	// previous one isn't a `defer` call.
	if _, ok := previousNode.(*ast.DeferStmt); !ok {
		w.CheckCuddling(stmt, cursor, 1)
	}
}

func (w *WSL) CheckBranch(stmt *ast.BranchStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckBreak]; !ok && stmt.Tok == token.BREAK {
		return
	}

	if _, ok := w.Config.Checks[CheckContinue]; !ok && stmt.Tok == token.CONTINUE {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	lastStmtInBlock := cursor.statements[len(cursor.statements)-1]
	if stmt == lastStmtInBlock && len(cursor.statements) <= 2 {
		return
	}

	w.addError(
		stmt.Pos(),
		stmt.Pos(),
		stmt.Pos(),
		MessageAddWhitespace,
	)
}

func (w *WSL) CheckDecl(stmt *ast.DeclStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckDecl]; !ok {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	w.CheckCuddling(stmt, cursor, 1)
}

func (w *WSL) CheckBlock(block *ast.BlockStmt) {
	w.CheckUnnecessaryBlockLeadingNewline(block)

	cursor := NewCursor(-1, block.List)
	for cursor.Next() {
		// fmt.Printf("%d: %T\n", cursor.currentIdx, cursor.Stmt())
		w.CheckStmt(cursor.Stmt(), cursor)
	}
}

func (w *WSL) CheckReturn(stmt *ast.ReturnStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckReturn]; !ok {
		return
	}

	// There's only a return statement.
	noStmts := cursor.Len()
	if noStmts <= 1 {
		return
	}

	// If the distance between the first statement and the return statement is
	// less than 3 LOC we're allowed to cuddle.
	firstStmts := cursor.Nth(0)
	if w.lineFor(stmt.End())-w.lineFor(firstStmts.Pos()) < 2 {
		return
	}

	w.addError(
		stmt.Pos(),
		stmt.Pos(),
		stmt.Pos(),
		MessageAddWhitespace,
	)
}

func (w *WSL) CheckAssign(stmt *ast.AssignStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckAssign]; ok {
		previousNode := cursor.PreviousNode()
		_, prevIsAssign := previousNode.(*ast.AssignStmt)
		_, prevIsDecl := previousNode.(*ast.DeclStmt)

		if w.numberOfStatementsAbove(cursor) > 0 && previousNode != nil && !prevIsAssign && !prevIsDecl {
			w.addError(
				stmt.Pos(),
				stmt.Pos(),
				stmt.Pos(),
				MessageAddWhitespace,
			)
		}
	}

	for _, expr := range stmt.Rhs {
		w.CheckExpr(expr, cursor)
	}
}

func (w *WSL) CheckStmt(stmt ast.Stmt, cursor *Cursor) {
	//nolint:gocritic // This is not commented out code, it's examples
	switch s := stmt.(type) {
	// if a {} else if b {} else {}
	case *ast.IfStmt:
		w.CheckIf(s, cursor)
	// for {} / for a; b; c {}
	case *ast.ForStmt:
		w.CheckFor(s, cursor)
	// for _, _ = range a {}
	case *ast.RangeStmt:
		w.CheckRange(s, cursor)
	// switch {} // switch a {}
	case *ast.SwitchStmt:
		w.CheckSwitch(s, cursor)
	// switch a.(type) {}
	case *ast.TypeSwitchStmt:
		w.CheckTypeSwitch(s, cursor)
	// return a
	case *ast.ReturnStmt:
		w.CheckReturn(s, cursor)
	// continue / break
	case *ast.BranchStmt:
		w.CheckBranch(s, cursor)
	// var a
	case *ast.DeclStmt:
		w.CheckDecl(s, cursor)
	// a := a
	case *ast.AssignStmt:
		w.CheckAssign(s, cursor)
	// a++ / a--
	case *ast.IncDecStmt:
	// defer func() {}
	case *ast.DeferStmt:
		w.CheckDefer(s, cursor)
	// go func() {}
	case *ast.GoStmt:
		w.CheckGo(s, cursor)
	case *ast.ExprStmt:
		w.CheckExprStmt(s, cursor)
	case *ast.CaseClause:
	case *ast.BlockStmt:
		w.CheckBlock(s)
	default:
		fmt.Printf("Not implemented stmt: %T\n", s)
	}
}

func (w *WSL) CheckExpr(expr ast.Expr, cursor *Cursor) {
	switch s := expr.(type) {
	// func() {}
	case *ast.FuncLit:
		w.CheckBlock(s.Body)
	// Call(args...)
	case *ast.CallExpr:
		for _, e := range s.Args {
			w.CheckExpr(e, cursor)
		}
	case *ast.BasicLit, *ast.CompositeLit, *ast.Ident, *ast.UnaryExpr:
	default:
		fmt.Printf("Not implemented expr: %T\n", s)
	}
}

// numberOfStatementsAbove will find out how many lines above the cursor's
// current statement there is without any newlines between.
func (w *WSL) numberOfStatementsAbove(cursor *Cursor) int {
	reset := cursor.Save()
	defer reset()

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
	if _, ok := w.Config.Checks[CheckLeadingWhitespace]; !ok {
		return
	}

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
	case *ast.ExprStmt:
		idents = append(idents, allIdents(n.X)...)
	case *ast.DeclStmt:
		idents = append(idents, allIdents(n.Decl)...)
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			idents = append(idents, allIdents(spec)...)
		}
	case *ast.GoStmt:
		idents = append(idents, allIdents(n.Call)...)
	case *ast.DeferStmt:
		idents = append(idents, allIdents(n.Call)...)
	case *ast.ValueSpec:
		for _, name := range n.Names {
			idents = append(idents, allIdents(name)...)
		}
	case *ast.AssignStmt:
		// TODO: For TypeSwitchStatements, this can be a false positive by
		// allowing shadowing and "tricking" usage;
		// var v any

		// notV := 1
		// switch notV := v.(type) {}
		//
		// This would trick wsl to see `notV` used in both type switch and on
		// line above - faulty(?)
		for _, lhs := range n.Lhs {
			idents = append(idents, allIdents(lhs)...)
		}

		// This must be here to see if a variable is used on the RHS, e.g.
		// a := 1
		// b = append(b, fmt.Sprintf("%s", a))
		for _, rhs := range n.Rhs {
			idents = append(idents, allIdents(rhs)...)
		}
	case *ast.IfStmt:
		idents = append(idents, allIdents(n.Cond)...)
		// TODO: idents = append(idents, allIdents(n.Else)...)
	case *ast.BinaryExpr:
		idents = append(idents, allIdents(n.X)...)
		idents = append(idents, allIdents(n.Y)...)
	case *ast.RangeStmt:
		idents = append(idents, allIdents(n.X)...)
	case *ast.SelectorExpr:
		idents = append(idents, allIdents(n.X)...)
	case *ast.UnaryExpr:
		idents = append(idents, allIdents(n.X)...)
	case *ast.ForStmt:
		idents = append(idents, allIdents(n.Init)...)
		idents = append(idents, allIdents(n.Cond)...)
		idents = append(idents, allIdents(n.Post)...)
	case *ast.SwitchStmt:
		idents = append(idents, allIdents(n.Init)...)
		idents = append(idents, allIdents(n.Tag)...)
	case *ast.TypeSwitchStmt:
		idents = append(idents, allIdents(n.Init)...)
		idents = append(idents, allIdents(n.Assign)...)
	case *ast.TypeAssertExpr:
		idents = append(idents, allIdents(n.X)...)
	case *ast.CallExpr:
		idents = append(idents, allIdents(n.Fun)...)
		for _, arg := range n.Args {
			idents = append(idents, allIdents(arg)...)
		}
	case *ast.CompositeLit:
		for _, elt := range n.Elts {
			idents = append(idents, allIdents(elt)...)
		}
	case *ast.BasicLit, *ast.FuncLit, *ast.IncDecStmt, *ast.BranchStmt,
		*ast.ArrayType:
	default:
		spew.Dump(node)
		fmt.Printf("missing ident detection for %T\n", node)
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

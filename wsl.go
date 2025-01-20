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

func (w *WSL) CheckCuddling(stmt ast.Node, cursor, blockCursor *Cursor, maxAllowedStatements int) {
	w.checkCuddlingWithDecl(stmt, cursor, blockCursor, maxAllowedStatements, true)
}

func (w *WSL) CheckCuddlingNoDecl(stmt ast.Node, cursor, blockCursor *Cursor, maxAllowedStatements int) {
	w.checkCuddlingWithDecl(stmt, cursor, blockCursor, maxAllowedStatements, false)
}

func (w *WSL) checkCuddlingWithDecl(
	stmt ast.Node,
	cursor *Cursor,
	blockCursor *Cursor,
	maxAllowedStatements int,
	declIsValid bool,
) {
	defer cursor.Save()()

	currentIdents := allIdents(cursor.Stmt())
	previousIdents := []*ast.Ident{}

	var previousNode ast.Node

	if cursor.Previous() {
		previousNode = cursor.Stmt()
		previousIdents = allIdents(previousNode)

		cursor.Next() // Move forward again
	}

	numStmtsAbove := w.numberOfStatementsAbove(cursor)

	// If we don't have any statements above, we only care about potential error
	// cuddling (for if statements) so check that.
	if numStmtsAbove == 0 {
		w.checkError(numStmtsAbove, stmt, previousNode, previousIdents, cursor)
		return
	}

	_, prevIsAssign := previousNode.(*ast.AssignStmt)
	_, prevIsDecl := previousNode.(*ast.DeclStmt)
	_, prevIsIncDec := previousNode.(*ast.IncDecStmt)
	_, currIsDefer := stmt.(*ast.DeferStmt)

	// Most of the time we allow cuddling with declarations (var) if it's just
	// one statement but not always so this can be disabled, e.g. for
	// delclarations themselves.
	if !declIsValid {
		prevIsDecl = false
	}

	// We're cuddled but not with an assign, declare or defer statement which is
	// never allowed.
	if !prevIsAssign && !prevIsDecl && !currIsDefer && !prevIsIncDec {
		w.addError(cursor.Stmt().Pos(), cursor.Stmt().Pos(), cursor.Stmt().Pos(), MessageAddWhitespace)
		return
	}

	// FEATURE: Allow identifier used anywhere in block (including recursive
	// blocks).
	if w.Config.AllowWholeBlock {
		anyIntersects := identsInMap(previousIdents, blockCursor.idents)
		if len(anyIntersects) > 0 {
			// We have matches, but too many statements above.
			if numStmtsAbove > maxAllowedStatements {
				w.addError(previousNode.Pos(), previousNode.Pos(), previousNode.Pos(), MessageAddWhitespace)
			}

			return
		}
	}

	// FEATURE: Allow identifiers used first in block. Configurable to allow
	// multiple levels.
	if !w.Config.AllowWholeBlock && w.Config.AllowFirstInBlock {
		for i := range w.Config.FirstInBlockMaxDepth {
			if i < len(blockCursor.firstIdents) {
				firstIntersect := identIntersection(
					previousIdents,
					blockCursor.firstIdents[i],
				)

				if len(firstIntersect) > 0 {
					// We have matches, but too many statements above.
					if numStmtsAbove > maxAllowedStatements {
						w.addError(previousNode.Pos(), previousNode.Pos(), previousNode.Pos(), MessageAddWhitespace)
					}

					return
				}
			}
		}
	}

	// We're cuddled but the line immediately above doesn't contain any
	// variables used in this statement.
	intersects := identIntersection(currentIdents, previousIdents)
	if len(intersects) == 0 {
		w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
		return
	}

	// Check the maximum number of allowed statements above, and if not
	// disabled (-1) check that the previous one intersects with the
	// current one.
	if maxAllowedStatements != -1 && numStmtsAbove > maxAllowedStatements {
		// Idents on the line above exist in the current condition so that
		// should remain cuddled.
		if len(intersects) > 0 {
			w.addError(previousNode.Pos(), previousNode.Pos(), previousNode.Pos(), MessageAddWhitespace)
		}
	}
}

func (w *WSL) CheckCuddlingWithoutIntersection(stmt ast.Node, cursor *Cursor) {
	previousNode := cursor.PreviousNode()

	previousAssign, prevIsAssign := previousNode.(*ast.AssignStmt)
	_, prevIsDecl := previousNode.(*ast.DeclStmt)
	_, prevIsIncDec := previousNode.(*ast.IncDecStmt)
	currAssign, currIsAssign := stmt.(*ast.AssignStmt)

	// Most of the time we allow cuddling with declarations (var) if it's just
	// one statement but not always so this can be disabled, e.g. for
	// declarations themselves.
	if _, ok := w.Config.Checks[CheckDecl]; ok {
		prevIsDecl = false
	}

	// If we enable exclusive assign checks we only allow new declarations or
	// new assignments together but not mix and match.
	// When this is enabled we also implicitly disable support to cuddle with
	// anything else.
	if _, ok := w.Config.Checks[CheckAssignExclusive]; ok {
		prevIsDecl = false
		prevIsIncDec = false

		if prevIsAssign && currIsAssign {
			prevIsAssign = previousAssign.Tok == currAssign.Tok
		}
	}

	prevIsValidType := previousNode == nil || prevIsAssign || prevIsDecl || prevIsIncDec

	// TODO: This is `allow-assign-and-call`/`AllowAssignAndCallCuddle`, should
	// we deprecate it?
	// ref: https://github.com/bombsimon/wsl/blob/52299dcd5c1c2a8baf77b4be4508937486d43656/wsl.go#L559-L563
	// 1. It's not actually checking call - just that we have intersections
	// 2. It's a bit too niche I think, either we support assign and call (or
	// whatever) or we don't.
	// 3. With the new check config, one could just disable checks for assign
	// and it would allow cuddling with anything.
	if !prevIsValidType && currIsAssign {
		currentIdents := allIdents(stmt)
		previousIdents := allIdents(previousNode)
		intersects := identIntersection(currentIdents, previousIdents)

		if len(intersects) > 0 {
			return
		}
	}

	if w.numberOfStatementsAbove(cursor) > 0 && !prevIsValidType {
		w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
	}
}

func (w *WSL) checkError(
	stmtsAbove int,
	ifStmt ast.Node,
	previousNode ast.Node,
	previousIdents []*ast.Ident,
	cursor *Cursor,
) {
	if _, ok := ifStmt.(*ast.IfStmt); !ok {
		return
	}

	if _, ok := w.Config.Checks[CheckErr]; !ok {
		return
	}

	if stmtsAbove > 0 || len(previousIdents) == 0 {
		return
	}

	if !slices.ContainsFunc(previousIdents, func(ident *ast.Ident) bool {
		return w.implementsErr(ident)
	}) {
		return
	}

	w.addError(ifStmt.Pos(), previousNode.End(), ifStmt.Pos(), MessageRemoveWhitespace)

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
		w.addError(ifStmt.Pos(), previousNode.Pos(), previousNode.Pos(), MessageAddWhitespace)
	}
}

func (w *WSL) CheckBlockAndExtend(
	node ast.Node,
	blockStmt *ast.BlockStmt,
	cursor *Cursor,
	check CheckType,
) {
	blockCursor := w.CheckBlock(blockStmt)

	if _, ok := w.Config.Checks[check]; ok {
		w.CheckCuddling(node, cursor, blockCursor, 1)
	}

	cursor.Extend(blockCursor)
}

func (w *WSL) CheckExprAndExtend(
	node ast.Node,
	expr ast.Expr,
	cursor *Cursor,
	predicate func(ast.Node) (int, bool),
	check CheckType,
) {
	maybeBlockCursor := w.CheckExpr(expr, cursor)

	if _, ok := w.Config.Checks[check]; ok {
		previousNode := cursor.PreviousNode()

		// We can cuddle any amount `go` statements so only check cuddling if the
		// previous one isn't a `go` call.
		// We don't even have to check if it's actually cuddled or just the previous
		// one because even if it's not but is a `go` statement it's valid.
		if n, ok := predicate(previousNode); ok {
			w.CheckCuddling(node, cursor, maybeBlockCursor, n)
		}
	}

	cursor.Extend(maybeBlockCursor)
}

func (w *WSL) CheckFunc(funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	w.CheckBlock(funcDecl.Body)
}

func (w *WSL) CheckIf(stmt *ast.IfStmt, cursor *Cursor) {
	// if
	blockCursor := w.CheckBlock(stmt.Body)

	switch v := stmt.Else.(type) {
	// else-if
	case *ast.IfStmt:
		// We use our new cursor created for this if-block to check nested if.
		w.CheckIf(v, blockCursor)
		blockCursor.Retain()

	// else
	case *ast.BlockStmt:
		// Merge the idents from the else branch.
		blockCursor.Merge(w.CheckBlock(v))
	}

	if _, ok := w.Config.Checks[CheckIf]; ok {
		w.CheckCuddling(stmt, cursor, blockCursor, 1)
	}

	cursor.Extend(blockCursor)
}

func (w *WSL) CheckFor(stmt *ast.ForStmt, cursor *Cursor) {
	w.CheckBlockAndExtend(stmt, stmt.Body, cursor, CheckFor)
}

func (w *WSL) CheckRange(stmt *ast.RangeStmt, cursor *Cursor) {
	w.CheckBlockAndExtend(stmt, stmt.Body, cursor, CheckRange)
}

func (w *WSL) CheckSwitch(stmt *ast.SwitchStmt, cursor *Cursor) {
	w.CheckBlockAndExtend(stmt, stmt.Body, cursor, CheckSwitch)
}

func (w *WSL) CheckTypeSwitch(stmt *ast.TypeSwitchStmt, cursor *Cursor) {
	w.CheckBlockAndExtend(stmt, stmt.Body, cursor, CheckTypeSwitch)
}

func (w *WSL) CheckExprStmt(stmt *ast.ExprStmt, cursor *Cursor) {
	w.CheckExprAndExtend(
		stmt,
		stmt.X,
		cursor,
		func(n ast.Node) (int, bool) {
			_, ok := n.(*ast.ExprStmt)
			return -1, !ok
		},
		CheckExpr,
	)
}

func (w *WSL) CheckGo(stmt *ast.GoStmt, cursor *Cursor) {
	w.CheckExprAndExtend(
		stmt,
		stmt.Call,
		cursor,
		func(n ast.Node) (int, bool) {
			_, ok := n.(*ast.GoStmt)
			return 1, !ok
		},
		CheckGo,
	)
}

func (w *WSL) CheckDefer(stmt *ast.DeferStmt, cursor *Cursor) {
	w.CheckExprAndExtend(
		stmt,
		stmt.Call,
		cursor,
		func(n ast.Node) (int, bool) {
			_, ok := n.(*ast.DeferStmt)
			return 1, !ok
		},
		CheckDefer,
	)
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

	// TODO: Should this be lines or statements?
	lastStmtInBlock := cursor.statements[len(cursor.statements)-1]
	if stmt == lastStmtInBlock && len(cursor.statements) <= 2 {
		return
	}

	w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
}

func (w *WSL) CheckDecl(stmt *ast.DeclStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckDecl]; !ok {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	w.CheckCuddlingNoDecl(stmt, cursor, &Cursor{}, 1)
}

func (w *WSL) CheckBlock(block *ast.BlockStmt) *Cursor {
	w.CheckBlockLeadingNewline(block)
	w.CheckTrailingNewline(block)

	return w.checkBody(block.List)
}

func (w *WSL) CheckCase(stmt *ast.CaseClause, cursor *Cursor) {
	w.CheckCaseLeadingNewline(stmt)

	if w.Config.CaseMaxLines != 0 {
		w.checkCaseTrailingNewline(stmt.Body, cursor)
	}

	cursor.Extend(w.checkBody(stmt.Body))
}

func (w *WSL) CheckComm(stmt *ast.CommClause, cursor *Cursor) {
	w.CheckCommLeadingNewline(stmt)

	if w.Config.CaseMaxLines != 0 {
		w.checkCaseTrailingNewline(stmt.Body, cursor)
	}

	cursor.Extend(w.checkBody(stmt.Body))
}

func (w *WSL) checkCaseTrailingNewline(body []ast.Stmt, cursor *Cursor) {
	if len(body) == 0 {
		return
	}

	defer cursor.Save()()

	if !cursor.Next() {
		return
	}

	var nextCase ast.Node
	switch n := cursor.Stmt().(type) {
	case *ast.CaseClause:
		nextCase = n
	case *ast.CommClause:
		nextCase = n
	default:
		return
	}

	firstStmt := body[0]
	lastStmt := body[len(body)-1]
	totalLines := w.lineFor(lastStmt.End()) - w.lineFor(firstStmt.Pos()) + 1

	// Not exceeding max lines to require newline.
	if totalLines <= w.Config.CaseMaxLines {
		return
	}

	// Next case is not immediately after the last statement so must be newline
	// already.
	if w.lineFor(nextCase.Pos()) > w.lineFor(lastStmt.End())+1 {
		return
	}

	w.addError(lastStmt.End(), nextCase.Pos(), nextCase.Pos(), MessageAddWhitespace)
}

func (w *WSL) CheckReturn(stmt *ast.ReturnStmt, cursor *Cursor) {
	for _, expr := range stmt.Results {
		w.CheckExpr(expr, cursor)
	}

	if _, ok := w.Config.Checks[CheckReturn]; !ok {
		return
	}

	// There's only a return statement.
	noStmts := cursor.Len()
	if noStmts <= 1 {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	// If the distance between the first statement and the return statement is
	// less than 3 LOC we're allowed to cuddle.
	firstStmts := cursor.Nth(0)
	if w.lineFor(stmt.End())-w.lineFor(firstStmts.Pos()) < w.Config.ReturnMaxLines {
		return
	}

	w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
}

func (w *WSL) CheckAssign(stmt *ast.AssignStmt, cursor *Cursor) {
	defer func() {
		for _, expr := range stmt.Rhs {
			w.CheckExpr(expr, cursor)
		}
	}()

	if _, ok := w.Config.Checks[CheckAssign]; !ok {
		return
	}

	w.CheckCuddlingWithoutIntersection(stmt, cursor)
	w.strictAppendCheck(stmt, cursor)
}

func (w *WSL) CheckIncDec(stmt *ast.IncDecStmt, cursor *Cursor) {
	defer w.CheckExpr(stmt.X, cursor)

	w.CheckCuddlingWithoutIntersection(stmt, cursor)
}

func (w *WSL) CheckLabel(stmt *ast.LabeledStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckLabel]; !ok {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
}

func (w *WSL) strictAppendCheck(stmt *ast.AssignStmt, cursor *Cursor) {
	if _, ok := w.Config.Checks[CheckAppend]; !ok {
		return
	}

	if w.numberOfStatementsAbove(cursor) == 0 {
		return
	}

	previousNode := cursor.PreviousNode()

	var appendNode *ast.CallExpr

	for _, expr := range stmt.Rhs {
		e, ok := expr.(*ast.CallExpr)
		if !ok {
			continue
		}

		if f, ok := e.Fun.(*ast.Ident); ok && f.Name == "append" {
			appendNode = e
			break
		}
	}

	if appendNode == nil {
		return
	}

	appendIdents := allIdents(appendNode)
	previousIdents := allIdents(previousNode)
	intersects := identIntersection(appendIdents, previousIdents)

	if len(intersects) == 0 {
		w.addError(stmt.Pos(), stmt.Pos(), stmt.Pos(), MessageAddWhitespace)
	}
}

func (w *WSL) checkBody(body []ast.Stmt) *Cursor {
	cursor := NewCursor(body)
	isFirst := true

	for cursor.Next() {
		w.CheckStmt(cursor.Stmt(), cursor)

		cursor.AddIdents(allIdents(cursor.Stmt()), isFirst)
		isFirst = false
	}

	return cursor
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
		w.CheckIncDec(s, cursor)
	// defer func() {}
	case *ast.DeferStmt:
		w.CheckDefer(s, cursor)
	// go func() {}
	case *ast.GoStmt:
		w.CheckGo(s, cursor)
	// e.g. someFn()
	case *ast.ExprStmt:
		w.CheckExprStmt(s, cursor)
	// case:
	case *ast.CaseClause:
		w.CheckCase(s, cursor)
	// case:
	case *ast.CommClause:
		w.CheckComm(s, cursor)
	// { }
	case *ast.BlockStmt:
		w.CheckBlock(s)
	// select { }
	case *ast.SelectStmt:
		w.CheckBlock(s.Body)
	// ch <- ...
	case *ast.SendStmt:
		// TODO: Check cuddling?
		w.CheckExpr(s.Value, cursor)
	// LABEL:
	case *ast.LabeledStmt:
		w.CheckLabel(s, cursor)
	default:
		fmt.Printf("Not implemented stmt: %T\n", s)
	}
}

func (w *WSL) CheckExpr(expr ast.Expr, cursor *Cursor) *Cursor {
	switch s := expr.(type) {
	// func() {}
	case *ast.FuncLit:
		return w.CheckBlock(s.Body)
	// Call(args...)
	case *ast.CallExpr:
		c := w.CheckExpr(s.Fun, cursor)
		for _, e := range s.Args {
			c.Extend(w.CheckExpr(e, cursor))
		}

		return c
	case *ast.IndexExpr:
		return w.CheckExpr(s.X, cursor)
	case *ast.StarExpr:
		return w.CheckExpr(s.X, cursor)
	case *ast.ArrayType,
		*ast.BasicLit,
		*ast.BinaryExpr,
		*ast.ChanType,
		*ast.CompositeLit,
		*ast.Ident,
		*ast.MapType,
		*ast.ParenExpr,
		*ast.SelectorExpr,
		*ast.SliceExpr,
		*ast.TypeAssertExpr,
		*ast.UnaryExpr:
		cursor.AddIdents(allIdents(s), false)
	default:
		fmt.Printf("Not implemented expr: %T\n", s)
	}

	return NewCursor([]ast.Stmt{})
}

// numberOfStatementsAbove will find out how many lines above the cursor's
// current statement there is without any newlines between.
func (w *WSL) numberOfStatementsAbove(cursor *Cursor) int {
	defer cursor.Save()()

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

func (w *WSL) CheckBlockLeadingNewline(body *ast.BlockStmt) {
	comments := ast.NewCommentMap(w.Fset, body, w.File.Comments)
	w.CheckLeadingNewline(body.Lbrace, body.List, comments)
}

func (w *WSL) CheckCaseLeadingNewline(caseClause *ast.CaseClause) {
	comments := ast.NewCommentMap(w.Fset, caseClause, w.File.Comments)
	w.CheckLeadingNewline(caseClause.Colon, caseClause.Body, comments)
}

func (w *WSL) CheckCommLeadingNewline(commClause *ast.CommClause) {
	comments := ast.NewCommentMap(w.Fset, commClause, w.File.Comments)
	w.CheckLeadingNewline(commClause.Colon, commClause.Body, comments)
}

func (w *WSL) CheckLeadingNewline(startPos token.Pos, body []ast.Stmt, comments ast.CommentMap) {
	if _, ok := w.Config.Checks[CheckLeadingWhitespace]; !ok {
		return
	}

	// No statements in the block, let's leave it as is.
	if len(body) == 0 {
		return
	}

	openLine := w.lineFor(startPos)
	openingPos := startPos + 1
	firstStmt := body[0].Pos()

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
				case openingPosLine == openLine &&
					commentStartLine == openLine+1:
					openingPos = comment.End()
				// The opening position has been updated, it's another comment.
				case openingPosLine != openLine:
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

func (w *WSL) CheckTrailingNewline(body *ast.BlockStmt) {
	if _, ok := w.Config.Checks[CheckTrailingWhitespace]; !ok {
		return
	}

	// No statements in the block, let's leave it as is.
	if len(body.List) == 0 {
		return
	}

	lastStmt := body.List[len(body.List)-1]

	// We don't want to force removal of the empty line for the last case since
	// it can be use for consistency and readability.
	if _, ok := lastStmt.(*ast.CaseClause); ok {
		return
	}

	closingPos := body.Rbrace
	lastStmtOrComment := lastStmt.End()

	// Empty label statements needs positional adjustment. #92
	if l, ok := lastStmt.(*ast.LabeledStmt); ok {
		if _, ok := l.Stmt.(*ast.EmptyStmt); ok {
			lastStmtOrComment = lastStmt.Pos()
		}
	}

	comments := ast.NewCommentMap(w.Fset, body, w.File.Comments)
	for _, commentGroup := range comments {
		for _, comment := range commentGroup {
			if comment.End() < closingPos && comment.Pos() > lastStmtOrComment {
				lastStmtOrComment = comment.End()
			}
		}
	}

	closingPosLine := w.Fset.PositionFor(closingPos, false).Line
	lastStmtLine := w.Fset.PositionFor(lastStmtOrComment, false).Line

	if closingPosLine != lastStmtLine+1 {
		w.addError(closingPos, lastStmtOrComment, closingPos, MessageRemoveWhitespace)
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

		for _, value := range n.Values {
			idents = append(idents, allIdents(value)...)
		}
	case *ast.AssignStmt:
		// NOTE: For TypeSwitchStatements, this can be a false positive by
		// allowing shadowing and "tricking" usage;
		//
		// 	   var v any
		//
		// 	   notV := 1
		// 	   switch notV := v.(type) {}
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
		idents = append(idents, allIdents(n.Init)...)
		idents = append(idents, allIdents(n.Cond)...)
		// TODO
		// idents = append(idents, allIdents(n.Else)...)
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
	case *ast.IncDecStmt:
		idents = allIdents(n.X)
	case *ast.CaseClause:
		for _, expr := range n.List {
			idents = append(idents, allIdents(expr)...)
		}
	case *ast.ReturnStmt:
		for _, r := range n.Results {
			idents = append(idents, allIdents(r)...)
		}
	case *ast.KeyValueExpr:
		idents = append(idents, allIdents(n.Key)...)
		idents = append(idents, allIdents(n.Value)...)
	case *ast.MapType:
		idents = append(idents, allIdents(n.Key)...)
		idents = append(idents, allIdents(n.Value)...)
	case *ast.StarExpr:
		idents = append(idents, allIdents(n.X)...)
	case *ast.IndexExpr:
		idents = append(idents, allIdents(n.X)...)
		idents = append(idents, allIdents(n.Index)...)
	case *ast.SliceExpr:
		idents = append(idents, allIdents(n.X)...)
		idents = append(idents, allIdents(n.Low)...)
		idents = append(idents, allIdents(n.High)...)
		idents = append(idents, allIdents(n.Max)...)
	case *ast.FieldList:
		for _, f := range n.List {
			idents = append(idents, allIdents(f)...)
		}
	case *ast.StructType:
		idents = append(idents, allIdents(n.Fields)...)
	case *ast.ParenExpr:
		idents = append(idents, allIdents(n.X)...)
	case *ast.SendStmt:
		idents = append(idents, allIdents(n.Chan)...)
		idents = append(idents, allIdents(n.Value)...)
	case *ast.ChanType:
		idents = append(idents, allIdents(n.Value)...)
	case *ast.CommClause:
		for _, s := range n.Body {
			idents = append(idents, allIdents(s)...)
		}
	case *ast.ArrayType,
		*ast.BasicLit,
		*ast.BlockStmt,
		*ast.BranchStmt,
		*ast.FuncLit,
		*ast.LabeledStmt,
		*ast.SelectStmt,
		*ast.TypeSpec:
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

func identsInMap(a []*ast.Ident, m map[string]struct{}) []*ast.Ident {
	intersects := []*ast.Ident{}

	for _, as := range a {
		if _, ok := m[as.Name]; ok {
			intersects = append(intersects, as)
		}
	}

	return intersects
}

package wsl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
)

type result struct {
	FileName   string
	LineNumber int
	Position   token.Position
	Reason     string
}

type processor struct {
	result   []result
	warnings []string
	fileSet  *token.FileSet
	file     *ast.File
}

// ProcessFiles takes a string slice with file names (full paths) and lints
// them.
func ProcessFiles(filenames []string) ([]result, []string) {
	p := NewProcessor()

	// Iterate over all files.
	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		p.process(filename, data)
	}

	return p.result, p.warnings
}

// NewProcessor will create a processor.
func NewProcessor() *processor {
	return &processor{
		result: []result{},
	}
}

func (p *processor) process(filename string, data []byte) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, data, parser.ParseComments)

	// If the file is not parsable let's add a syntax error and move on.
	if err != nil {
		p.result = append(p.result, result{
			FileName:   filename,
			LineNumber: 0,
			Reason:     fmt.Sprintf("invalid syntax, file cannot be linted (%s)", err.Error()),
		})

		return
	}

	p.fileSet = fileSet
	p.file = file

	for _, d := range p.file.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			p.parseBlockBody(v.Body)
		case *ast.GenDecl:
			// `go fmt` will handle proper spacing for GenDecl such as imports,
			// constants etc.
		default:
			p.addWarning("type not implemented", d.Pos(), v)
		}
	}
}

// parseBlockBody will parse any kind of block statements such as switch cases
// and if statements. A list of Result is returned.
func (p *processor) parseBlockBody(block *ast.BlockStmt) {
	// Start by finding leading and trailing whitespaces.
	p.findLeadingAndTrailingWhitespaces(block, nil)

	// Parse the block body contents.
	p.parseBlockStatements(block.List)
}

// parseBlockStatements will parse all the statements found in the body of a
// node. A list of Result is returned.
func (p *processor) parseBlockStatements(statements []ast.Stmt) {
	for i, stmt := range statements {
		// Start by checking if the statement has a body (probably if-statement,
		// a range, switch case or similar. Whenever a body is found we start by
		// parsing it before moving on in the AST.
		statementBody := reflect.Indirect(reflect.ValueOf(stmt)).FieldByName("Body")

		// Some cases allow cuddling depending on the first statement in a body
		// of a block or case. If possible extract the first statement.
		var firstBodyStatement ast.Node

		if statementBody.IsValid() {
			switch statementBodyContent := statementBody.Interface().(type) {
			case *ast.BlockStmt:
				if len(statementBodyContent.List) > 0 {
					firstBodyStatement = statementBodyContent.List[0]
				}

				p.parseBlockBody(statementBodyContent)
			case []ast.Stmt:
				// The Body field for an *ast.CaseClause is of type []ast.Stmt.
				// We must check leading and trailing whitespaces and then pass
				// the statements to parseBlockStatements to parse it's content.
				var nextStatement ast.Node

				// Check if there's more statements (potential casess) after the
				// current one.
				if len(statements)-1 > i {
					nextStatement = statements[i+1]
				}

				p.findLeadingAndTrailingWhitespaces(stmt, nextStatement)
				p.parseBlockStatements(statementBodyContent)
			default:
				p.addWarning(
					"body statement type not implemented ",
					stmt.Pos(), statementBodyContent,
				)

				continue
			}
		}

		// First statement, nothing to do.
		if i == 0 {
			continue
		}

		previousStatement := statements[i-1]

		// If the last statement didn't end one line above the current statement
		// we know we're not cuddled so just move on.
		if p.fileSet.Position(previousStatement.End()).Line != p.fileSet.Position(stmt.Pos()).Line-1 {
			continue
		}

		// We know we're cuddled, extract assigned variables on the line above
		// which is the only thing we allow cuddling with. If the assignment is
		// made over multiple lines we should not allow cuddling.
		assignedOnLineAbove := []string{}

		// Ensure previous line is not a multi line assignment and if not get
		// all assigned variables.
		if p.fileSet.Position(previousStatement.Pos()).Line == p.fileSet.Position(stmt.Pos()).Line-1 {
			assignedOnLineAbove = p.findAssignments(previousStatement)
		}

		// We could potentially have a block which require us to check the first
		// argument before ruling out an allowed cuddle.
		assignedFirstInBlock := []string{}

		if firstBodyStatement != nil {
			assignedFirstInBlock = p.findAssignments(firstBodyStatement)
		}

		// Get usages on the line above, i.e. x.Foo() -> x
		usagesOnLineAbove := p.findExpressionVariables(previousStatement)

		switch t := stmt.(type) {
		case *ast.IfStmt:
			// Check if we're cuddled with something that's not an assignment.
			if len(assignedOnLineAbove) == 0 {
				p.addError(t.Pos(), "if statements should only be cuddled with assignments")

				continue
			}

			// Check if there are more than one statement before us.
			if i >= 2 {
				statementBeforePreviousStatement := statements[i-2]

				if p.fileSet.Position(statementBeforePreviousStatement.End()).Line == p.fileSet.Position(previousStatement.Pos()).Line-1 {
					p.addError(t.Pos(), "only one cuddle assignment allowed before if statement")

					continue
				}
			}

			// Get all variables used in the condition for the if statement.
			usedInStatement := p.findConditionVariables(t.Cond)

			if !atLeastOneInListsMatch(usedInStatement, assignedOnLineAbove) {
				// No variables in if statement was assigned in on the line
				// above, check if what's cuddled above is used on the first
				// line in the block.
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					// If the variable above isn't used first in the block this
					// is just too tight code.
					p.addError(t.Pos(), "if statements should only be cuddled with assignments used in the if statement itself")
				}

				continue
			}
		case *ast.ReturnStmt:
			// If we're the last statement, check if there's no more than two
			// lines from the starting statement and the end of this statement.
			// This is to support short return functions such as:
			// func (t *Typ) X() {
			//     t.X = true
			//     return t
			// }
			if i == len(statements)-1 && i == 1 {
				currentStatementEndsAtLine := p.fileSet.Position(stmt.End()).Line
				previousStatementStartsAtLine := p.fileSet.Position(previousStatement.Pos()).Line

				if currentStatementEndsAtLine-previousStatementStartsAtLine <= 2 {

					continue
				}
			}

			p.addError(t.Pos(), "return statements should never be cuddled")
		case *ast.AssignStmt:
			if _, ok := previousStatement.(*ast.AssignStmt); ok {
				continue
			}

			p.addError(t.Pos(), "assignments should only be cuddled with other assignments")
		case *ast.DeclStmt:
			p.addError(t.Pos(), "declarations should never be cuddled")
		case *ast.ExprStmt:
			switch previousStatement.(type) {
			case *ast.DeclStmt, *ast.ReturnStmt:
				p.addError(t.Pos(), "expressions should not be cuddled with declarations or returns")
			}
		case *ast.RangeStmt:
			rangesOverValues := p.findConditionVariables(t.X)

			// The same rule applies for ranges as for if statements, see
			// comments regarding variable usages on the line before or as the
			// first line in the block for details.
			if !atLeastOneInListsMatch(rangesOverValues, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					p.addError(t.Pos(), "ranges should only be cuddled with assignments used in the iteration")
				}
			}
		case *ast.DeferStmt:
			deferExpressionUsages := p.findCallExpressionVariables(t.Call)

			if !atLeastOneInListsMatch(deferExpressionUsages, usagesOnLineAbove) {
				p.addError(t.Pos(), "defer statements should only be cuddled with expressions on same variable")
			}
		case *ast.ForStmt:
			// No condition, just a for loop - should not be cuddled.
			if t.Cond == nil {
				p.addError(t.Pos(), "for statement without condition should never be cuddled")

				continue
			}

			rangesOverValues := p.findConditionVariables(t.Cond)

			// The same rule applies for ranges as for if statements, see
			// comments regarding variable usages on the line before or as the
			// first line in the block for details.
			if !atLeastOneInListsMatch(rangesOverValues, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					p.addError(t.Pos(), "for statements should only be cuddled with assignments used in the iteration")
				}
			}
		case *ast.GoStmt:
			if c, ok := t.Call.Fun.(*ast.Ident); ok {
				if atLeastOneInListsMatch([]string{c.Name}, assignedOnLineAbove) {
					continue
				}
			}

			p.addError(t.Pos(), "go statements can only invoke functions assigned on line above")
		case *ast.BranchStmt:
			// TODO: This is a continue or break within loop
			// Should probably be handled just like return
		case *ast.SwitchStmt:
			// TODO: Handle switch
		case *ast.TypeSwitchStmt:
			// TODO: Handle type switch
		case *ast.CommClause:
			// TODO: Handle default cases in select
		case *ast.CaseClause:
			// Case clauses will be checked by not allowing leading ot trailing
			// whitespaces within the block. There's nothing in the case itself
			// that may be cuddled.
		default:
			p.addWarning("stmt type not implemented", t.Pos(), t)
		}
	}
}

// findExpressionVariables will find used identifiers for a node, i.e.
// mu.Unlock() -> mu
func (p *processor) findExpressionVariables(node ast.Node) []string {
	astExpression, ok := node.(*ast.ExprStmt)
	if !ok {
		return []string{}
	}

	return p.findCallExpressionVariables(astExpression.X)
}

func (p *processor) findCallExpressionVariables(node ast.Node) []string {
	switch x := node.(type) {
	case *ast.CallExpr:
		exp, ok := x.Fun.(*ast.SelectorExpr)
		if !ok {
			return []string{}
		}

		return p.findConditionVariables(exp.X)
	default:
		p.addWarning("expression variable not implemented", x.Pos(), x)
	}

	return []string{}
}

// findAssignments will find all assignments for a given node. If the node isn't
// an assignment statement, an empty list will be return. If it however is an
// assignment statement, all identifiers on the left hand side will be returned.
func (p *processor) findAssignments(node ast.Node) []string {
	assignments := []string{}

	astIdentifier, ok := node.(*ast.AssignStmt)
	if !ok {
		return assignments
	}

	for _, identifier := range astIdentifier.Lhs {
		switch x := identifier.(type) {
		case *ast.Ident:
			assignments = append(assignments, x.Name)
		case *ast.SelectorExpr:
			// No new variables on left side.
		case *ast.IndexExpr:
			assignments = append(p.findConditionVariables(x.X))
		default:
			p.addWarning("assignment type not implemented", x.Pos(), x)
		}
	}

	return assignments
}

func (p *processor) findConditionVariables(cond ast.Node) []string {
	switch v := cond.(type) {
	case *ast.Ident:
		return []string{v.Name}
	case *ast.BinaryExpr:
		return append(
			p.findConditionVariables(v.X),
			p.findConditionVariables(v.Y)...,
		)
	case *ast.UnaryExpr:
		return p.findConditionVariables(v.X)
	case *ast.IndexExpr:
		return p.findConditionVariables(v.X)
	case *ast.CallExpr:
		args := []string{}

		for _, x := range v.Args {
			if callExprArg, ok := x.(*ast.Ident); ok {
				args = append(args, callExprArg.Name)
			}
		}

		return args
	case *ast.BasicLit:
		// Basic literal, nothing to check.
	case *ast.SelectorExpr:
		return p.findConditionVariables(v.X)
	default:
		p.addWarning("condition type not implemented", v.Pos(), v)
	}

	return []string{}
}

func atLeastOneInListsMatch(listOne, listTwo []string) bool {
	for _, lOne := range listOne {
		for _, lTwo := range listTwo {
			if lOne == lTwo {
				return true
			}
		}
	}

	return false
}

// findLeadingAndTrailingWhitespaces will find leading and trailing whitespaces
// in a node. The method takes comments in consideration which will make the
// parser more gentle. A list of Result is returned.
func (p *processor) findLeadingAndTrailingWhitespaces(stmt, nextStatement ast.Node) {
	var (
		allowedLinesBeforeFirstStatement = 1
		commentMap                       = ast.NewCommentMap(p.fileSet, stmt, p.file.Comments)
		firstStatement                   ast.Node
		lastStatement                    ast.Node
	)

	switch t := stmt.(type) {
	case *ast.BlockStmt:
		if len(t.List) < 1 {
			return
		}

		firstStatement = t.List[0]
		lastStatement = t.List[len(t.List)-1]
	case *ast.CaseClause:
		if len(t.Body) < 1 {
			return
		}

		firstStatement = t.Body[0]
		lastStatement = t.Body[len(t.Body)-1]
	case *ast.CommClause:
		if len(t.Body) < 1 {
			return
		}

		firstStatement = t.Body[0]
		lastStatement = t.Body[len(t.Body)-1]
	default:
		p.addWarning("whitespace node type not implemented ", stmt.Pos(), stmt)
	}

	if c, ok := commentMap[firstStatement]; ok {
		commentsBefore := c[0]

		// Ensure that the comment we are parsing occurs BEFORE the statement.
		if commentsBefore.Pos() < firstStatement.Pos() {
			allowedLinesBeforeFirstStatement += len(commentsBefore.List)
		}
	}

	switch t := stmt.(type) {
	case *ast.BlockStmt:
		if p.fileSet.Position(firstStatement.Pos()).Line != p.fileSet.Position(t.Lbrace).Line+allowedLinesBeforeFirstStatement {
			p.addErrorOffset(
				firstStatement.Pos(),
				-1,
				"block should not start with a whitespace",
			)
		}

		if p.fileSet.Position(lastStatement.End()).Line != p.fileSet.Position(t.Rbrace).Line-1 {
			p.addErrorOffset(
				lastStatement.Pos(),
				1,
				"block should not end with a whitespace (or comment)",
			)
		}
	case *ast.CaseClause:
		if p.fileSet.Position(firstStatement.Pos()).Line != p.fileSet.Position(t.Colon).Line+allowedLinesBeforeFirstStatement {
			p.addErrorOffset(
				lastStatement.Pos(),
				-1,
				"case block should not start with a whitespace",
			)
		}

		if nextStatement != nil {
			if nextStatementCase, ok := nextStatement.(*ast.CaseClause); ok {
				if p.fileSet.Position(lastStatement.End()).Line != p.fileSet.Position(nextStatementCase.Colon).Line-1 {
					p.addErrorOffset(
						lastStatement.Pos(),
						1,
						"case block should not end with a whitespace (or comment)",
					)
				}
			}
		}
	}
}

// Add an error for the file and line number for the current token.Pos with the
// given reason.
func (p *processor) addError(pos token.Pos, reason string) {
	p.addErrorOffset(pos, 0, reason)
}

// Add an error for the file for the current token.Pos with the given offset and
// reason. The offset will be added to the token.Pos line.
func (p *processor) addErrorOffset(pos token.Pos, offset int, reason string) {
	position := p.fileSet.Position(pos)

	p.result = append(p.result, result{
		FileName:   position.Filename,
		LineNumber: position.Line + offset,
		Position:   position,
		Reason:     reason,
	})
}

func (p *processor) addWarning(w string, pos token.Pos, t interface{}) {
	position := p.fileSet.Position(pos)

	p.warnings = append(p.warnings,
		fmt.Sprintf("%s:%d: %s (%T)", position.Filename, position.Line, w, t),
	)
}

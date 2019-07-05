package wsl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
)

// Result represents the result of one linted row which was not accepted by the
// white space linter.
type Result struct {
	FileName   string
	LineNumber int
	Pos        token.Position
	Reason     string
}

// ProcessFiles takes a string slice with file names (full paths) and lints
// them.
func ProcessFiles(files []string) []Result {
	var result []Result

	for _, file := range files {
		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		result = append(result, ProcessFile(file, fileData)...)
	}

	return result
}

// ProcessFile will process a single file by reading the sorce and parsing it to
// a token.FileSet which holds the full ast (abstract syntax tree) for the file.
// A list of Result is returned.
func ProcessFile(fileName string, fileData []byte) []Result {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, fileData, parser.ParseComments)
	if err != nil {
		return []Result{
			{
				FileName:   file.Name.Name,
				LineNumber: 0,
				Reason:     fmt.Sprintf("invalid syntax, file cannot be linted (%s)", err.Error()),
			},
		}
	}

	return ProcessAST(fset, file)
}

// ProcessAST will process the file (or AST of file) via the *token.FileSet and
// *ast.File
func ProcessAST(fset *token.FileSet, file *ast.File) []Result {
	result := []Result{}

	for _, d := range file.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			result = append(result, parseBlockBody(fset, file.Comments, v.Body)...)
		case *ast.GenDecl:
			// `go fmt` will handle proper spacing for GenDecl such as imports,
			// constants etcetera.
		default:
			fmt.Printf("%s:%d: type not implemented (%T)\n", file.Name.Name, fset.Position(d.Pos()).Line, v)
		}
	}

	return result
}

// parseBlockBody will parse any kind of block statements such as switch cases
// and if statements. A list of Result is returned.
func parseBlockBody(fset *token.FileSet, comments []*ast.CommentGroup, block *ast.BlockStmt) []Result {
	result := []Result{}

	result = append(result, findLeadingAndTrailingWhitespaces(fset, block, nil, comments)...)
	result = append(result, parseBlockStatements(fset, comments, block.List)...)

	return result
}

// parseBlockStatements will parse all the statements found in the body of a
// node. A list of Result is returned.
func parseBlockStatements(fset *token.FileSet, comments []*ast.CommentGroup, statements []ast.Stmt) []Result {
	result := []Result{}

	for i, stmt := range statements {
		// Start by checking if the statement has a body (probably if-statement,
		// a range, switch case or similar. Whenever a body is found we start by
		// parsing it before moving on in the AST.
		statementBody := reflect.Indirect(reflect.ValueOf(stmt)).FieldByName("Body")

		// Some cases allow cuddling depending on thefirst statement in a body
		// of a block or case. If possible extract the first staement.
		var firstBodyStatement ast.Node

		if statementBody.IsValid() {
			switch statementBodyContent := statementBody.Interface().(type) {
			case *ast.BlockStmt:
				if len(statementBodyContent.List) > 0 {
					firstBodyStatement = statementBodyContent.List[0]
				}

				result = append(result, parseBlockBody(fset, comments, statementBodyContent)...)
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

				result = append(result, findLeadingAndTrailingWhitespaces(fset, stmt, nextStatement, comments)...)
				result = append(result, parseBlockStatements(fset, comments, statementBodyContent)...)
			default:
				fmt.Printf("%s:%d: body statement type not implemented (%T)\n", fset.File(stmt.Pos()).Name(), fset.Position(stmt.Pos()).Line, statementBodyContent)
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
		if fset.Position(previousStatement.End()).Line != fset.Position(stmt.Pos()).Line-1 {
			continue
		}

		// We know we're cuddled, extract assigned variables on the line above
		// which is the only thing we allow cuddling with. If the assignment is
		// made over multiple lines we should not allow cuddling.
		assignedOnLineAbove := []string{}

		// Ensure previous line is not a multi line assignment and if not get
		// all assigned variables.
		if fset.Position(previousStatement.Pos()).Line == fset.Position(stmt.Pos()).Line-1 {
			assignedOnLineAbove = findAssignments(previousStatement, fset)
		}

		// We could potentially have a block which require us to check the first
		// argument before ruling out an allowed cuddle.
		assignedFirstInBlock := []string{}

		if firstBodyStatement != nil {
			assignedFirstInBlock = findAssignments(firstBodyStatement, fset)
		}

		// Get usages on the line above, i.e. x.Foo() -> x
		usagesOnLineAbove := findExpressionVariables(previousStatement, fset)

		switch t := stmt.(type) {
		case *ast.IfStmt:
			// Check if we're cuddled with something that's not an assigment.
			if len(assignedOnLineAbove) == 0 {
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Pos:        fset.Position(t.Pos()),
					Reason:     "if statements should only be cuddled with assigments",
				})

				continue
			}

			// Check if there are more than one statement before us.
			if i >= 2 {
				statementBeforePreviousStatement := statements[i-2]

				if fset.Position(statementBeforePreviousStatement.End()).Line == fset.Position(previousStatement.Pos()).Line-1 {
					result = append(result, Result{
						FileName:   fset.Position(t.Pos()).Filename,
						LineNumber: fset.Position(t.Pos()).Line,
						Pos:        fset.Position(t.Pos()),
						Reason:     "only one cuddle assignment allowed before if statement",
					})

					continue
				}
			}

			// Get all variables used in the condition for the if statement.
			usedInStatement := findConditionVariables(t.Cond, fset)

			if !atLeastOneInListsMatch(usedInStatement, assignedOnLineAbove) {
				// No variables in if statement was assigned in on the line
				// above, check if what's cuddled above is used on the first
				// line in the block.
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					// If the variable above isn't used first in the block this
					// is just too tight code.
					result = append(result, Result{
						FileName:   fset.Position(t.Pos()).Filename,
						LineNumber: fset.Position(t.Pos()).Line,
						Pos:        fset.Position(t.Pos()),
						Reason:     "if statements should only be cuddled with assigments used in the if statement itself",
					})
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
				currentStatementEndsAtLine := fset.Position(stmt.End()).Line
				previousStatementStartsAtLine := fset.Position(previousStatement.Pos()).Line

				if currentStatementEndsAtLine-previousStatementStartsAtLine <= 2 {
					continue
				}
			}

			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Pos:        fset.Position(t.Pos()),
				Reason:     "return statements should never be cuddled",
			})
		case *ast.AssignStmt:
			if _, ok := previousStatement.(*ast.AssignStmt); ok {
				continue
			}

			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Pos:        fset.Position(t.Pos()),
				Reason:     "assigments should only be cuddled with other assigments",
			})
		case *ast.DeclStmt:
			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Pos:        fset.Position(t.Pos()),
				Reason:     "declarations should never be cuddled",
			})
		case *ast.ExprStmt:
			switch previousStatement.(type) {
			case *ast.DeclStmt, *ast.ReturnStmt:
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Pos:        fset.Position(t.Pos()),
					Reason:     "expressions should not be cuddled with declarations or returns",
				})
			}
		case *ast.RangeStmt:
			rangesOverValues := findConditionVariables(t.X, fset)

			// The same rule applies for ranges as for if statements, see
			// comments regarding variable usages on the line before or as the
			// first line in the block for details.
			if !atLeastOneInListsMatch(rangesOverValues, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					result = append(result, Result{
						FileName:   fset.Position(t.Pos()).Filename,
						LineNumber: fset.Position(t.Pos()).Line,
						Pos:        fset.Position(t.Pos()),
						Reason:     "ranges should only be cuddled with assignments used in the iteration",
					})
				}
			}
		case *ast.DeferStmt:
			deferExpressionUsages := findCallExpressionVariables(t.Call, fset)

			if !atLeastOneInListsMatch(deferExpressionUsages, usagesOnLineAbove) {
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Pos:        fset.Position(t.Pos()),
					Reason:     "defer statements should only be cuddled with expressions on same variable",
				})
			}
		case *ast.ForStmt:
			// No condition, just a for loop - should not be cuddled.
			if t.Cond == nil {
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Pos:        fset.Position(t.Pos()),
					Reason:     "for statement without condition should never be cuddled",
				})

				continue
			}

			rangesOverValues := findConditionVariables(t.Cond, fset)

			// The same rule applies for ranges as for if statements, see
			// comments regarding variable usages on the line before or as the
			// first line in the block for details.
			if !atLeastOneInListsMatch(rangesOverValues, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					result = append(result, Result{
						FileName:   fset.Position(t.Pos()).Filename,
						LineNumber: fset.Position(t.Pos()).Line,
						Pos:        fset.Position(t.Pos()),
						Reason:     "for statements should only be cuddled with assignments used in the iteration",
					})
				}
			}
		case *ast.GoStmt:
			if c, ok := t.Call.Fun.(*ast.Ident); ok {
				if atLeastOneInListsMatch([]string{c.Name}, assignedOnLineAbove) {
					continue
				}
			}

			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Pos:        fset.Position(t.Pos()),
				Reason:     "go statements can only invoke functions assigned on line above",
			})
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
			fmt.Printf("%s:%d: stmt type not implemented (%T)\n", fset.File(t.Pos()).Name(), fset.Position(t.Pos()).Line, t)
		}
	}

	return result
}

// findExpressionVariables will find used identifiers for a node, i.e.
// mu.Unlock() -> mu
func findExpressionVariables(node ast.Node, fset *token.FileSet) []string {
	astExpression, ok := node.(*ast.ExprStmt)
	if !ok {
		return []string{}
	}

	return findCallExpressionVariables(astExpression.X, fset)
}

func findCallExpressionVariables(node ast.Node, fset *token.FileSet) []string {
	switch x := node.(type) {
	case *ast.CallExpr:
		exp, ok := x.Fun.(*ast.SelectorExpr)
		if !ok {
			return []string{}
		}

		return findConditionVariables(exp.X, fset)
	default:
		fmt.Printf("%s:%d: expression variable not implemented (%T)\n", fset.File(x.Pos()).Name(), fset.Position(x.Pos()).Line, x)
	}

	return []string{}
}

// findassignments will find all asignmenets for a given node. If the node isn't
// an assignment statement, an empty list will be return. If it however is an
// assignment statement, all identifiers on the left hand side will be returned.
func findAssignments(node ast.Node, fset *token.FileSet) []string {
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
			assignments = findConditionVariables(x.X, fset)
		default:
			fmt.Printf("%s:%d: assignment type not implemented (%T)\n", fset.File(x.Pos()).Name(), fset.Position(x.Pos()).Line, x)
		}
	}

	return assignments
}

func findConditionVariables(cond ast.Node, fset *token.FileSet) []string {
	switch v := cond.(type) {
	case *ast.Ident:
		return []string{v.Name}
	case *ast.BinaryExpr:
		return append(
			findConditionVariables(v.X, fset),
			findConditionVariables(v.Y, fset)...,
		)
	case *ast.UnaryExpr:
		return findConditionVariables(v.X, fset)
	case *ast.IndexExpr:
		return findConditionVariables(v.X, fset)
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
		return findConditionVariables(v.X, fset)
	default:
		fmt.Printf("%s:%d: condition type not implemented (%T)\n", fset.File(v.Pos()).Name(), fset.Position(v.Pos()).Line, v)
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
func findLeadingAndTrailingWhitespaces(fset *token.FileSet, stmt, nextStatement ast.Node, comments []*ast.CommentGroup) []Result {
	var (
		result                           = []Result{}
		allowedLinesBeforeFirstStatement = 1
		commentMap                       = ast.NewCommentMap(fset, stmt, comments)
		firstStatement                   ast.Node
		lastStatement                    ast.Node
	)

	switch t := stmt.(type) {
	case *ast.BlockStmt:
		if len(t.List) < 1 {
			return result
		}

		firstStatement = t.List[0]
		lastStatement = t.List[len(t.List)-1]
	case *ast.CaseClause:
		if len(t.Body) < 1 {
			return result
		}

		firstStatement = t.Body[0]
		lastStatement = t.Body[len(t.Body)-1]
	case *ast.CommClause:
		if len(t.Body) < 1 {
			return result
		}

		firstStatement = t.Body[0]
		lastStatement = t.Body[len(t.Body)-1]
	default:
		fmt.Printf("%s:%d: whitespace node type not implemented (%T)\n", fset.File(stmt.Pos()).Name(), fset.Position(stmt.Pos()).Line, stmt)
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
		if fset.Position(firstStatement.Pos()).Line != fset.Position(t.Lbrace).Line+allowedLinesBeforeFirstStatement {
			result = append(result, Result{
				FileName:   fset.Position(firstStatement.Pos()).Filename,
				LineNumber: fset.Position(firstStatement.Pos()).Line - 1,
				Pos:        fset.Position(firstStatement.Pos()),
				Reason:     "block should not start with a whitespace",
			})
		}

		if fset.Position(lastStatement.End()).Line != fset.Position(t.Rbrace).Line-1 {
			result = append(result, Result{
				FileName:   fset.Position(lastStatement.Pos()).Filename,
				LineNumber: fset.Position(lastStatement.Pos()).Line + 1,
				Pos:        fset.Position(lastStatement.Pos()),
				Reason:     "block should not end with a whitespace (or comment)",
			})
		}
	case *ast.CaseClause:
		if fset.Position(firstStatement.Pos()).Line != fset.Position(t.Colon).Line+allowedLinesBeforeFirstStatement {
			result = append(result, Result{
				FileName:   fset.Position(firstStatement.Pos()).Filename,
				LineNumber: fset.Position(firstStatement.Pos()).Line - 1,
				Pos:        fset.Position(firstStatement.Pos()),
				Reason:     "case block should not start with a whitespace",
			})
		}

		if nextStatement != nil {
			if nextStatementCase, ok := nextStatement.(*ast.CaseClause); ok {
				if fset.Position(lastStatement.End()).Line != fset.Position(nextStatementCase.Colon).Line-1 {
					result = append(result, Result{
						FileName:   fset.Position(lastStatement.Pos()).Filename,
						LineNumber: fset.Position(lastStatement.Pos()).Line + 1,
						Pos:        fset.Position(lastStatement.Pos()),
						Reason:     "case block should not end with a whitespace (or comment)",
					})
				}
			}
		}
	}

	return result
}

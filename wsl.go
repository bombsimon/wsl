package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
)

// Result represnets the result of one linted row which was not accepted by the
// white space linter.
type Result struct {
	FileName   string
	LineNumber int
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
	result := []Result{}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, fileData, parser.ParseComments)
	if err != nil {
		return []Result{
			{
				FileName:   file.Name.Name,
				LineNumber: 0,
				Reason:     "invalid syntax, file cannot be linted",
			},
		}
	}

	for _, d := range file.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			result = append(result, parseBlock(fset, file.Comments, v.Body)...)
		case *ast.GenDecl:
			// `go fmt` will handle proper spacing for GenDecl such as imports,
			// constants etcetera.
		default:
			fmt.Printf("%s:%d: type not implemented (%T)\n", file.Name.Name, fset.Position(d.Pos()).Line, v)
		}
	}

	return result
}

// parseBlock will parse any kind of block statements such as switch cases and
// if statements. A list of Result is returned.
func parseBlock(fset *token.FileSet, comments []*ast.CommentGroup, block *ast.BlockStmt) []Result {
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
		bodyV := reflect.Indirect(reflect.ValueOf(stmt)).FieldByName("Body")

		if bodyV.IsValid() {
			switch bodyT := bodyV.Interface().(type) {
			case *ast.BlockStmt:
				result = append(result, parseBlock(fset, comments, bodyT)...)
			case *ast.CaseClause:
				// TODO: Must handle starting and traingling whitespace
				result = append(result, parseBlockStatements(fset, comments, bodyT.Body)...)
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
				result = append(result, parseBlockStatements(fset, comments, bodyT)...)
			default:
				fmt.Printf("%s:%d: body statement type not implemented (%T)\n", fset.File(stmt.Pos()).Name(), fset.Position(stmt.Pos()).Line, bodyT)
				continue
			}
		}

		// First statement, nothing to do.
		if i == 0 {
			continue
		}

		previousStatement := statements[i-1]

		// The last statement did not end at the line above (must be a
		// whitespace, nothing to check)
		if fset.Position(previousStatement.End()).Line != fset.Position(stmt.Pos()).Line-1 {
			continue
		}

		switch t := stmt.(type) {
		case *ast.IfStmt:
			// Check if we're cuddled with something that's not an assigment.
			lastAssignment, ok := previousStatement.(*ast.AssignStmt)
			if !ok {
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Reason:     "if statements can only be cuddled with assigments",
				})

				continue
			}

			switch t.Cond.(type) {
			case *ast.UnaryExpr, *ast.BinaryExpr, *ast.Ident:
				// Expected expressions found in if condition, handle below.
			case *ast.SelectorExpr:
				// TODO: Never cuddle?!
				continue
			default:
				fmt.Printf("%s:%d: condition type not implemented (%T)\n", fset.File(t.Pos()).Name(), fset.Position(t.Pos()).Line, t)
				continue
			}

			foundAssigmentOnLineAbove := false

			for _, astIdentifier := range lastAssignment.Lhs {
				previousAssigmentIdentifier, ok := astIdentifier.(*ast.Ident)
				if !ok {
					// Ignore assignments which isn't idents. This will
					// eventually end up with foundAssigmentOnLineAbove being
					// false which is handled.
					continue
				}

				var expressionInIfStatement *ast.Ident

				switch x := t.Cond.(type) {
				case *ast.Ident:
					expressionInIfStatement = x
				case *ast.BinaryExpr:
					switch expressionType := x.X.(type) {
					case *ast.Ident:
						expressionInIfStatement = expressionType
					case *ast.CallExpr:
						// TODO: Recursive check args to function X and/or Y
						continue
					}
				case *ast.UnaryExpr:
					expressionInIfStatement = x.X.(*ast.Ident)
				default:
					panic(fmt.Sprintf("Not an expression at %d?! %T\n", x, fset.Position(t.Pos()).Line))
				}

				if previousAssigmentIdentifier.Name == expressionInIfStatement.Name {
					foundAssigmentOnLineAbove = true
					break
				}
			}

			if !foundAssigmentOnLineAbove {
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Reason:     "if statements can only be cuddled with assigments used in the if statement itself",
				})

				continue
			}
		case *ast.ReturnStmt:
			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Reason:     "return statements can never be cuddled",
			})
		case *ast.AssignStmt:
			if _, ok := previousStatement.(*ast.AssignStmt); ok {
				continue
			}

			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Reason:     "assigments can only be cuddled with other assigments",
			})
		case *ast.DeclStmt:
			result = append(result, Result{
				FileName:   fset.Position(t.Pos()).Filename,
				LineNumber: fset.Position(t.Pos()).Line,
				Reason:     "declarations can never be cuddled",
			})
		case *ast.ExprStmt:
			switch previousStatement.(type) {
			case *ast.DeclStmt, *ast.ReturnStmt:
				result = append(result, Result{
					FileName:   fset.Position(t.Pos()).Filename,
					LineNumber: fset.Position(t.Pos()).Line,
					Reason:     "expressions can not be cuddled with decarations or returns",
				})
			}
		case *ast.RangeStmt:
		// TODO: Handle for range...
		case *ast.BranchStmt:
			// TODO: What is this?
		case *ast.CaseClause:
			// Already handeled
		default:
			fmt.Printf("%s:%d: stmt type not implemented (%T)\n", fset.File(t.Pos()).Name(), fset.Position(t.Pos()).Line, t)
		}
	}

	return result
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
				Reason:     "block should not start with a whitespace",
			})
		}

		if fset.Position(lastStatement.End()).Line != fset.Position(t.Rbrace).Line-1 {
			result = append(result, Result{
				FileName:   fset.Position(lastStatement.Pos()).Filename,
				LineNumber: fset.Position(lastStatement.Pos()).Line + 1,
				Reason:     "block should not end with a whitespace",
			})
		}
	case *ast.CaseClause:
		if fset.Position(firstStatement.Pos()).Line != fset.Position(t.Colon).Line+allowedLinesBeforeFirstStatement {
			result = append(result, Result{
				FileName:   fset.Position(firstStatement.Pos()).Filename,
				LineNumber: fset.Position(firstStatement.Pos()).Line - 1,
				Reason:     "case block should not start with a whitespace",
			})
		}

		if nextStatement != nil {
			if nextStatementCase, ok := nextStatement.(*ast.CaseClause); ok {
				if fset.Position(lastStatement.End()).Line != fset.Position(nextStatementCase.Colon).Line-1 {
					result = append(result, Result{
						FileName:   fset.Position(lastStatement.Pos()).Filename,
						LineNumber: fset.Position(lastStatement.Pos()).Line + 1,
						Reason:     "case block should not end with a whitespace",
					})
				}
			}
		}
	}

	return result
}

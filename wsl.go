package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"

	"github.com/davecgh/go-spew/spew"
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
		result = append(result, ProcessAST(file)...)
	}

	return result
}

func ProcessAST(fileName string) []Result {
	result := []Result{}

	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

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
			// Ignore global declarations for now...
		default:
			fmt.Printf("not implemented type at line %d: %T\n", fset.Position(v.Pos()).Line, v)
		}
	}

	return result
}

func parseBlock(fset *token.FileSet, comments []*ast.CommentGroup, block *ast.BlockStmt) []Result {
	var (
		commentMap = ast.NewCommentMap(fset, block, comments)
		result     = []Result{}
	)

	// Nothing to parse, empty block???
	if len(block.List) < 1 {
		return []Result{}
	}

	firstStatement := block.List[0]
	lastStatement := block.List[len(block.List)-1]

	allowedLinesBeforeFirstStatement := 1
	if c, ok := commentMap[firstStatement]; ok {
		commentsBefore := c[0]

		// Ensure that the comment we are parsing occurs BEFORE the statement.
		if commentsBefore.Pos() < firstStatement.Pos() {
			allowedLinesBeforeFirstStatement += len(commentsBefore.List)
		}
	}

	if fset.Position(firstStatement.Pos()).Line != fset.Position(block.Lbrace).Line+allowedLinesBeforeFirstStatement {
		result = append(result, Result{
			FileName:   fset.Position(firstStatement.Pos()).Filename,
			LineNumber: fset.Position(firstStatement.Pos()).Line - 1,
			Reason:     "block should not start with a whitespace",
		})
	}

	if fset.Position(lastStatement.End()).Line != fset.Position(block.Rbrace).Line-1 {
		result = append(result, Result{
			FileName:   fset.Position(lastStatement.Pos()).Filename,
			LineNumber: fset.Position(lastStatement.Pos()).Line + 1,
			Reason:     "block should not end with a whitespace",
		})
	}

	result = append(result, parseBlockList(fset, comments, block.List)...)

	return result
}

func parseBlockList(fset *token.FileSet, comments []*ast.CommentGroup, statements []ast.Stmt) []Result {
	result := []Result{}

	for i, stmt := range statements {
		bodyV := reflect.Indirect(reflect.ValueOf(stmt)).FieldByName("Body")

		if bodyV.IsValid() {
			switch bodyT := bodyV.Interface().(type) {
			case *ast.BlockStmt:
				result = append(result, parseBlock(fset, comments, bodyT)...)
			case *ast.CaseClause:
				// TODO: Must handle starting and traingling whitespace
				result = append(result, parseBlockList(fset, comments, bodyT.Body)...)
			case []ast.Stmt:
			default:
				spew.Dump(bodyT)
				panic("what is this?!")
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

			switch expression := t.Cond.(type) {
			case *ast.BinaryExpr, *ast.Ident:
				// Expected expressions found in if condition, handle below.
			case *ast.SelectorExpr:
				// TODO: Never cuddler?!
				continue
			default:
				fmt.Printf("Unhandled expression IN if statement at line %d type: %T\n", fset.Position(t.Pos()).Line, expression)
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
					Reason:     "if statements can only be cuddled with assigments used in the if statement itself (only other assigments above)",
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
			// TODO: Cannot be cuddled with DeclStmt or ReturnStmts
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
		default:
			fmt.Printf("Unimplementedstatement at line %d:: %T\n", fset.Position(t.Pos()).Line, t)
		}
	}

	return result
}

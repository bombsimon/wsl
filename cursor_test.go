package wsl

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavestate(t *testing.T) {
	cursor := NewCursor(0, []ast.Stmt{
		&ast.AssignStmt{},
		&ast.AssignStmt{},
		&ast.AssignStmt{},
		&ast.AssignStmt{},
		&ast.AssignStmt{},
		&ast.AssignStmt{},
		&ast.AssignStmt{},
	})

	cursor.Save()
	cursor.Next()
	cursor.Next()

	func() {
		cursor.Save()
		defer cursor.Reset()

		cursor.Next()
		cursor.Next()

		func() {
			cursor.Save()
			defer cursor.Reset()

			cursor.Next()
			cursor.Next()

			assert.Equal(t, 6, cursor.currentIdx)
		}()

		assert.Equal(t, 4, cursor.currentIdx)
	}()

	assert.Equal(t, 2, cursor.currentIdx)

	cursor.Reset()
	assert.Equal(t, 0, cursor.currentIdx)
}

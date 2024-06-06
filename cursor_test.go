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

	reset := cursor.Save()
	cursor.Next()
	cursor.Next()

	func() {
		reset := cursor.Save()
		defer reset()

		cursor.Next()
		cursor.Next()

		func() {
			reset := cursor.Save()
			defer reset()

			cursor.Next()
			cursor.Next()

			assert.Equal(t, 6, cursor.currentIdx)
		}()

		assert.Equal(t, 4, cursor.currentIdx)
	}()

	assert.Equal(t, 2, cursor.currentIdx)

	reset()
	assert.Equal(t, 0, cursor.currentIdx)
}

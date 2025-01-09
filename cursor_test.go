package wsl

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavestate(t *testing.T) {
	cursor := NewCursor([]ast.Stmt{
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

			assert.Equal(t, 5, cursor.currentIdx)
		}()

		assert.Equal(t, 3, cursor.currentIdx)
	}()

	assert.Equal(t, 1, cursor.currentIdx)

	reset()
	assert.Equal(t, -1, cursor.currentIdx)
}

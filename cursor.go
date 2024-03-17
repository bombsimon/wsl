package wsl

import (
	"errors"
	"go/ast"
)

var ErrCursourOutObFounds = errors.New("out of bounds")

type Cursor struct {
	currentIdx int
	savedIdx   int
	statements []ast.Stmt
}

func NewCursor(i int, statements []ast.Stmt) *Cursor {
	return &Cursor{
		currentIdx: i,
		statements: statements,
	}
}

func (c *Cursor) Next() bool {
	if c.currentIdx >= len(c.statements)-1 {
		return false
	}

	c.currentIdx++

	return true
}

func (c *Cursor) Previous() bool {
	if c.currentIdx == 0 {
		return false
	}

	c.currentIdx--

	return true
}

func (c *Cursor) Stmt() ast.Stmt {
	return c.statements[c.currentIdx]
}

func (c *Cursor) Save() {
	c.savedIdx = c.currentIdx
}

func (c *Cursor) Reset() {
	c.currentIdx = c.savedIdx
}

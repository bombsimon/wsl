package wsl

import (
	"errors"
	"go/ast"
)

var ErrCursourOutObFounds = errors.New("out of bounds")

type Cursor struct {
	currentIdx int
	statements []ast.Stmt
	saves      []int
}

func NewCursor(i int, statements []ast.Stmt) *Cursor {
	return &Cursor{
		currentIdx: i,
		statements: statements,
		saves:      []int{},
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
	if c.currentIdx <= 0 {
		return false
	}

	c.currentIdx--

	return true
}

func (c *Cursor) PeekNext() bool {
	c.Save()
	defer c.Reset()

	return c.Next()
}

func (c *Cursor) PeekPrevious() bool {
	c.Save()
	defer c.Reset()

	return c.Previous()
}

func (c *Cursor) Stmt() ast.Stmt {
	return c.statements[c.currentIdx]
}

func (c *Cursor) Save() {
	c.saves = append(c.saves, c.currentIdx)
}

func (c *Cursor) Reset() {
	if len(c.saves) == 0 {
		return
	}

	idx := c.saves[len(c.saves)-1]
	c.saves = c.saves[:len(c.saves)-1]
	c.currentIdx = idx
}

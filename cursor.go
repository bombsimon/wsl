package wsl

import (
	"errors"
	"go/ast"
)

var ErrCursourOutObFounds = errors.New("out of bounds")

type Cursor struct {
	currentIdx int
	statements []ast.Stmt
	idents     map[string]struct{}
}

func NewCursor(statements []ast.Stmt) *Cursor {
	return &Cursor{
		currentIdx: -1,
		statements: statements,
		idents:     make(map[string]struct{}),
	}
}

func (c *Cursor) AddIdents(idents []*ast.Ident) {
	for _, i := range idents {
		c.idents[i.Name] = struct{}{}
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

func (c *Cursor) PreviousNode() ast.Node {
	defer c.Save()()

	var previousNode ast.Node
	if c.Previous() {
		previousNode = c.Stmt()
	}

	return previousNode
}

func (c *Cursor) PeekNext() bool {
	reset := c.Save()
	defer reset()

	return c.Next()
}

func (c *Cursor) PeekPrevious() bool {
	reset := c.Save()
	defer reset()

	return c.Previous()
}

func (c *Cursor) Stmt() ast.Stmt {
	return c.statements[c.currentIdx]
}

func (c *Cursor) Save() func() {
	idx := c.currentIdx
	return func() {
		c.currentIdx = idx
	}
}

func (c *Cursor) Len() int {
	return len(c.statements)
}

func (c *Cursor) Nth(n int) ast.Stmt {
	return c.statements[n]
}

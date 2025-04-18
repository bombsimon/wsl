package wsl

import (
	"errors"
	"go/ast"
)

var ErrCursourOutObFounds = errors.New("out of bounds")

// Cursor holds a list of statements and a pointer to where in the list we are.
// Each block gets a new cursor and can be used to check previous or coming
// statements.
//
// A cursor also keeps a set of all identifiers seen in the block and child
// blocks and a separate list of string slices of first identifier in blocks.
// Each element represents the depth of where the identifiers was found.
//
// In this example:
//
//	    func Fn() {
//		       if true {
//		           a, b := 1, 2
//		           if false {
//		               c := 3
//		               d := 4
//		           }
//		       }
//		   }
//
// idents will be a map containing [a, b, c, d] and first identifiers will be
// [[a, b], [c]].
type Cursor struct {
	currentIdx  int
	statements  []ast.Stmt
	idents      map[string]struct{}
	firstIdents [][]*ast.Ident
}

// NewCursor creates a new cursor with a given list of statements.
func NewCursor(statements []ast.Stmt) *Cursor {
	return &Cursor{
		currentIdx: -1,
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
	defer c.Save()()
	return c.Next()
}

func (c *Cursor) PeekPrevious() bool {
	defer c.Save()()
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

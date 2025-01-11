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
// [[a, b], [c]]
type Cursor struct {
	currentIdx  int
	statements  []ast.Stmt
	idents      map[string]struct{}
	firstIdents [][]*ast.Ident
}

// NewCursor creates a new cursor with a given list of statements.
func NewCursor(statements []ast.Stmt) *Cursor {
	return &Cursor{
		currentIdx:  -1,
		statements:  statements,
		idents:      make(map[string]struct{}),
		firstIdents: [][]*ast.Ident{},
	}
}

// Extend will extend a cursor with another. This is used when returning from a
// check of a nested block to include all scopes recursively.
func (c *Cursor) Extend(otherCursor *Cursor) {
	for i := range otherCursor.idents {
		c.idents[i] = struct{}{}
	}

	c.firstIdents = append(c.firstIdents, otherCursor.firstIdents...)
}

// Merge a cursor with another. This is used when checking chained blocks such
// as if-else-chains.
func (c *Cursor) Merge(otherCursor *Cursor) {
	if len(otherCursor.firstIdents) > 0 {
		otherFirstIdents := otherCursor.firstIdents[len(otherCursor.firstIdents)-1]

		if len(c.firstIdents) == 0 {
			c.firstIdents = append(c.firstIdents, otherFirstIdents)
		} else {
			for _, v := range otherFirstIdents {
				c.firstIdents[0] = append(c.firstIdents[0], v)
			}
		}
	}

	for v := range otherCursor.idents {
		c.idents[v] = struct{}{}
	}
}

// Retain will retain the first idents from a Cursor putting them all at the
// first index. This is used when there are multiple blocks at the same level
// with different first identifiers such as if-else.
func (c *Cursor) Retain() {
	if len(c.firstIdents) == 0 {
		return
	}

	for _, firsts := range c.firstIdents[1:] {
		c.firstIdents[0] = append(c.firstIdents[0], firsts...)
	}

	c.firstIdents = c.firstIdents[:1]
}

// AddIdents adds all idents in the list to the set of identifiers. If it's the
// first statement for the cursor first identifiers will be extended.
func (c *Cursor) AddIdents(idents []*ast.Ident, addFirst bool) {
	for _, i := range idents {
		c.idents[i.Name] = struct{}{}
	}

	if addFirst {
		c.firstIdents = append(c.firstIdents, idents)
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

package database

import (
	"fmt"
	"runtime"
)

// Tree is the representation of a single parsed string.
type Tree struct {
	Root      *ListNode
	Lex       *Lexer
	peekCount int
	token     [3]Item
}

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be
// given the specified name. If an error is encountered, parsing stops and an
// empty map is returned with the error.
// next returns the next token.
func Parse(text string) {

	t := New()
	t.Parse(text)

}

func (t *Tree) next() Item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.Lex.nextItem()
	}
	return t.token[t.peekCount]
}

// peek returns but does not consume the next token.
func (t *Tree) peek() Item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.Lex.nextItem()
	return t.token[0]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// nextNonSpace returns the next non-space token.
func (t *Tree) nextNonSpace() (token Item) {
	for {
		token = t.next()
		if token.Typ != ItemSpace {
			break
		}
	}
	return token
}

// peekNonSpace returns but does not consume the next non-space token.
func (t *Tree) peekNonSpace() (token Item) {
	for {
		token = t.next()
		if token.Typ != ItemSpace {
			break
		}
	}
	t.backup()
	return token
}

// New allocates a new parse tree.
func New() *Tree {
	return &Tree{}
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string) {
	format = fmt.Sprintf("%d: %s", t.token[0].Line, format)
	panic(fmt.Errorf(format))
}

// error terminates processing.
func (t *Tree) error(err error) {
	t.errorf(fmt.Sprintf("%s", err))
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.Lex.drain()
		}
		*errp = e.(error)
	}
}

// Parse parses the template definition string to construct a representation of
// the template for execution. If either action delimiter string is empty, the
// default ("{{" or "}}") is used. Embedded template definitions are added to
// the treeSet map.
func (t *Tree) Parse(text string) (err error) {
	defer t.recover(&err)
	t.Lex = Lex(text)
	s := t.peek().Pos
	//	fmt.Printf("%T", s)
	t.Root = t.newList(s)

	t.parse()
	return nil
}

// parse is the top-level parser for a template, essentially the same
// as ItemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) parse() {

	for t.peek().Typ != ItemEOF {

		switch n := t.textOrAction(); n.Type() {
		case nodeEnd:
			t.errorf(fmt.Sprintf("unexpected %s", n))
		default:
			t.Root.append(n)
			fmt.Println(n)
		}

	}
}

// textOrAction:
//	text | action
func (t *Tree) textOrAction() Node {
	switch token := t.nextNonSpace(); {
	case token.Typ == ItemText:
		return t.newText(token.Pos, token.Val)
	case token.Typ == ItemLeftDelimiter || token.Typ == ItemRightDelimiter || token.Typ == ItemAsterisk:
		return t.newText(token.Pos, token.Val)
	default:
		t.errorf(fmt.Sprintf("unexpected token: %s", token))
	}
	return nil
}

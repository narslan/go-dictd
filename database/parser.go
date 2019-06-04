package database

import (
	"fmt"
)

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be
// given the specified name. If an error is encountered, parsing stops and an
// empty map is returned with the error.
// next returns the next token.
func Parse(text string) {

	t := NewTree()
	t.Parse(text)

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
	t.Root = t.NewList(s)

	t.parse()
	return nil
}

// parse is the top-level parser for a template, essentially the same
// as ItemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) parse() {

	for t.peek().Typ != ItemEOF {

		switch n := t.textOrAction(); n.Type() {
		case NodeEnd:
			t.errorf(fmt.Sprintf("unexpected %s", n))
		default:
			t.Root.Append(n)
			fmt.Println(n)
		}

	}
}

// textOrAction:
//	text | action
func (t *Tree) textOrAction() Node {
	/*	switch token := t.nextNonSpace(); {
		case token.Typ == ItemText:
			return t.NewText(token.Pos, token.Val)
		case token.Typ == ItemLeftDelimiter || token.Typ == ItemRightDelimiter || token.Typ == ItemAsterisk:
			return t.NewText(token.Pos, token.Val)
		default:
			t.errorf(fmt.Sprintf("unexpected token: %s", token))
		}*/
	return nil
}

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

	t.Root = t.NewList(0)
	t.Lex = Lex(text)
	t.Parse()

}

// parse is the top-level parser for a template, essentially the same
// as ItemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) Parse() {

	// Root Selector contains all the roots.
	rs := make([]*ListNode, 0)
	depth := 0
	rs = append(rs, t.Root)
Loop:
	for {

		switch n := t.textOrAction(); n.Type() {
		case NodeList:
			fmt.Printf("New list node is created at depth %d\n", depth)
			t := n.(*ListNode)
			rs = append(rs, t)
			depth++
		case NodeEnd:
			depth--
			fmt.Printf("The list node is closed at depth %d\n", depth)
		case NodeString:
			fmt.Printf("item is added at depth %d %s\n", depth, n)
			rs[depth].Append(n)
		default:
			break Loop
		}

	}
}

// textOrAction:
//	text | action
func (t *Tree) textOrAction() Node {
	switch token := t.nextNonSpace(); {
	case token.Typ == ItemText:
		return t.NewString(token.Val, token.Pos)
	case token.Typ == ItemLeftDelimiter:
		return t.action()
	case token.Typ == ItemRightDelimiter:
		return t.end()
	default:
		return t.endParse()
	}
	return nil
}

func (t *Tree) action() (n Node) {
	token := t.peek()
	return t.NewList(token.Pos)
}

func (t *Tree) end() (n Node) {
	token := t.peek()
	return t.NewEnd(token.Pos)
}

func (t *Tree) endParse() (n Node) {
	token := t.peek()
	return t.NewBreak(token.Pos)
}

package database

import (
	_ "fmt"
)

func Parse(text string) []*ListNode {

	t := NewTree()

	t.Lex = Lex(text)

	t.Root = t.NewList(t.peek().Pos)

	// Root Selector contains all the roots.
	rs := make([]*ListNode, 0)
	depth := 0
	rs = append(rs, t.Root)

LOOP:
	for t.peek().Typ != ItemEOF {

		switch n := t.Action(); n.Type() {
		case NodeList:
			//	fmt.Printf("New list node is created at depth %d\n", depth)
			t := n.(*ListNode)
			rs = append(rs, t)
			depth++
		case NodeClose:
			depth--
		//	fmt.Printf("The list node is closed at depth %d\n", depth)
		case NodeString:
			//fmt.Printf("item is added at depth %d %s\n", depth, n)
			rs[depth].Append(n)
		case NodeEnd:
			break LOOP
		}

	}
	return rs
}

// textOrAction:
//	text | action
func (t *Tree) Action() Node {
	switch token := t.nextNonSpace(); {
	case token.Typ == ItemText:
		return t.NewString(token.Val, token.Pos)
	case token.Typ == ItemLeftDelimiter:
		return t.NewList(token.Pos)
	case token.Typ == ItemRightDelimiter:
		return t.NewClose(token.Pos)
	default:
		return t.NewEnd(token.Pos)
	}
	return nil
}

type DictdFile struct {
	Name  string
	Data  string
	Index string
}

func ParseDatabases(text string) []*DictdFile {

	t := NewTree()

	t.Lex = Lex(text)

	t.Root = t.NewList(t.peek().Pos)

	dictc := make([]*DictdFile, 0)
LOOP:
	for t.peek().Typ != ItemEOF {
		token := t.nextNonSpace()
		switch {
		case token.Val == "database":
			file := &DictdFile{}
			dictc = append(dictc, file)
			dictname := t.nextNonSpace()
			file.Name = dictname.Val
			del := t.nextNonSpace()
			if del.Typ == ItemLeftDelimiter {
				datad := t.nextNonSpace()
				if datad.Val == "data" {
					datadv := t.nextNonSpace()
					file.Data = datadv.Val

					dataiv := t.nextNonSpace()
					if dataiv.Val == "index" {
						indexv := t.nextNonSpace()
						file.Index = indexv.Val
					}

				}

			} else {
				t.errorf("{ expected")
			}

			//			fmt.Printf("%v\n", token.Val)
		case token.Typ == ItemEOF:
			break LOOP
		}
	}
	return dictc
}

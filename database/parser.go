package database

import (
	"io/ioutil"
	"log"
)

// Parse parses dictd configuration file format in a tree. As it is not used at the moment,
// it is useless. But it will be used as a boilerplate for further
// development.

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

// Action determines what type of token flows in.
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

// ReadConfig reads a dictd configuration file and
// returns the available dictionaries.
// TODO: Handling parse errors
func ReadConfig(filePath string) []*DictdDatabase {

	c, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	s := string(c)

	t := NewTree()

	t.Lex = Lex(s)

	t.Root = t.NewList(t.peek().Pos)

	dictc := make([]*DictdDatabase, 0)

LOOP:
	for t.peek().Typ != ItemEOF {
		token := t.nextNonSpace()
		switch {
		case token.Val == "database":
			//file := &DictdDatabase{}
			//dictc = append(dictc, file)

			name := t.nextNonSpace().Val
			var indexPath, dbPath string

			delimiter := t.nextNonSpace()

			if delimiter.Typ == ItemLeftDelimiter {
				//"data" token, consume it.
				datad := t.nextNonSpace()
				// Find data and consume.
				if datad.Val == "data" {
					dbPath = t.nextNonSpace().Val
					// Find index.
					dataiv := t.nextNonSpace()
					if dataiv.Val == "index" {
						indexPath = t.nextNonSpace().Val
					}

				}
				d, err := NewDictdDatabase(name, indexPath, dbPath)
				if err != nil {
					log.Fatal(err)
				}
				dictc = append(dictc, d)
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

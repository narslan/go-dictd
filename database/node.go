package database

import (
	"bytes"
	"fmt"
	"runtime"
	_ "strconv"
	_ "strings"
)

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

type Tree struct {
	Root      *ListNode
	peekCount int
	Lex       *Lexer
	token     [3]Item // three-token lookahead for parser.

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

// NewTree allocates and returns a Tree.
//

func NewTree() *Tree {
	return &Tree{}
}

type Node interface {
	Type() NodeType
	String() string
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeString = iota // A string constant.
	NodeList          // A List of Nodes.
	NodeMap
	NodeEnd
	NodeBreak
)

type ListNode struct {
	NodeType
	Pos
	Nodes []Node
}

func (t *Tree) NewList(pos Pos) *ListNode {
	return &ListNode{NodeType: NodeList, Pos: pos}
}

func (l *ListNode) Append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) String() string {
	b := new(bytes.Buffer)
	for _, n := range l.Nodes {
		fmt.Fprint(b, n)
	}
	return b.String()
}

type StringNode struct {
	NodeType
	Pos
	Text string
}

func (t *Tree) NewString(text string, pos Pos) *StringNode {
	return &StringNode{NodeType: NodeString, Text: text, Pos: pos}
}

func (s *StringNode) String() string {
	return "+" + s.Text
}

type MapNode struct {
	NodeType
	Pos
	Value map[string]string
}

func (t *Tree) NewMap(key, value string, pos Pos) *MapNode {
	m := make(map[string]string)
	m[key] = value
	return &MapNode{NodeType: NodeMap, Value: m, Pos: pos}
}

func (m MapNode) String() (s string) {

	for k, v := range m.Value {
		s = "Key:" + k + "Value:" + v
	}
	return s
}

// endNode represents an {{end}} action.
// It does not appear in the final parse tree.
type BreakNode struct {
	NodeType
	Pos
}

func (t *Tree) NewBreak(pos Pos) *BreakNode {
	return &BreakNode{NodeType: NodeBreak, Pos: pos}
}

func (e *BreakNode) String() string {
	return "end"
}

type EndNode struct {
	NodeType
	Pos
}

func (t *Tree) NewEnd(pos Pos) *EndNode {
	return &EndNode{NodeType: NodeEnd, Pos: pos}
}

func (e *EndNode) String() string {
	return "end"
}

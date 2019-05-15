package database

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Pos represents a byte position in the original input text from which
// this template was parsed.
//type Pos int

type ItemType int

// Item represents a token or text string returned from the scanner.
type Item struct {
	Typ  ItemType // The type of this Item.
	Pos  Pos      // The starting position, in bytes, of this Item in the input string.
	Val  string   // The value of this Item.
	Line int      // The line number at the start of this Item.
}

const eof = -1

const (
	ItemError          ItemType = iota // error occurred; value is text of error
	ItemAsterisk                       // printable ascii character; grab bag for *
	ItemLeftDelimiter                  // { character
	ItemRightDelimiter                 // } character
	ItemEOF                            //
	ItemSpace                          // run of spaces separating arguments
	ItemText                           // plain text
)

func (i Item) String() string {

	switch {
	case i.Typ == ItemEOF:
		return "EOF\n"
	case i.Typ == ItemError:
		return i.Val
	case i.Typ == ItemSpace:
		return ""
	default:
		return fmt.Sprintf("Keyword :  %d/%d %s \n", i.Line, i.Pos, i.Val)
	}
	return ""

}

// stateFn represents the state of the scanner as a function that returns the
// next state.

type stateFn func(*Lexer) stateFn

// Lexer holds the state of the scanner.
type Lexer struct {
	Input     string    // the string being scanned
	Pos       Pos       // current position in the input
	Start     Pos       // start position of this Item
	Width     Pos       // width of last rune read from input
	Items     chan Item // channel of scanned Items
	Line      int       // 1+number of newlines seen
	StartLine int       // start line of this Item
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.Pos) >= len(l.Input) {
		l.Width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width = Pos(w)
	l.Pos += l.Width
	if r == '\n' {
		l.Line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.Pos -= l.Width
	// Correct newline count.
	if l.Width == 1 && l.Input[l.Pos] == '\n' {
		l.Line--
	}
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Items <- Item{ItemError, l.Start, fmt.Sprintf(format, args...), l.StartLine}
	return nil
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.Line += strings.Count(l.Input[l.Start:l.Pos], "\n")
	l.Start = l.Pos
	l.StartLine = l.Line
}

// emit passes an Item back to the client.
func (l *Lexer) emit(t ItemType) {
	l.Items <- Item{t, l.Start, l.Input[l.Start:l.Pos], l.StartLine}
	l.Start = l.Pos
	l.StartLine = l.Line
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) nextItem() Item {
	return <-l.Items
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) drain() {
	for range l.Items {
	}
}

// lexText scans until an opening action delimiter, "{{".
func lexText(l *Lexer) stateFn {

	switch r := l.next(); {
	case r == eof:
		l.emit(ItemEOF)
		return nil
	case isSpace(r) || isEndOfLine(r):
		return lexSpace
	case r == '#':
		return lexComment
	case r == '{':
		l.emit(ItemLeftDelimiter)
	case r == '}':
		l.emit(ItemRightDelimiter)
	case r == '*':
		l.emit(ItemAsterisk)
	case r == '\'' || r == '"':
		return lexQuote
	case isAlphaNumeric(r):
		return lexWord
	}

	return lexText
}

// lexQuote scans until the end of a quote
func lexQuote(l *Lexer) stateFn {

	// This ignores the quote character.
	l.ignore()
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"', '\'':
			break Loop
		}
	}

	// This ignores the second quote character.
	l.backup()
	l.emit(ItemText)
	// We put the ignored character back into reading process.
	l.next()
	return lexText

}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {

	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(ItemSpace)
	return lexText
}

// lexComment ignores the input

func lexComment(l *Lexer) stateFn {
	l.Pos += Pos(1)
	i := strings.Index(l.Input[l.Pos:], "\n")
	l.Pos += Pos(i)
	l.ignore()
	return lexText
}

// lexWord scans a word.
func lexWord(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			//word := l.Input[l.Start:l.Pos]
			//	fmt.Printf("Catched >%s<  ", word)
			//_, ok := preDefined[word]
			//switch {
			//case ok:
			//l.emit(ItemKeyword)
			//default:
			l.emit(ItemText)
			//}

			break Loop
		}
	}
	return lexText
}

// lex creates a new scanner for the input string.
func Lex(input string) *Lexer {

	l := &Lexer{
		Input:     input,
		Items:     make(chan Item),
		Line:      1,
		StartLine: 1,
	}
	go l.run()
	return l
}

// run runs the state machine for the Lexer.
func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.Items)
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '-' || r == '.' || r == '/' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

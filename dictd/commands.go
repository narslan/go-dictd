package dictd

import (
	"errors"
	"strings"
	"unicode"
)

/* Command is the encapsulation for a user's request of the Server. */

type Command struct {
	Command string
	Params  []string
}

/* Parse an incoming line, and return a `dict.Command` suitable for
 * passing to internal (or external) handlers.
 * *
 *https://groups.google.com/forum/#!topic/golang-nuts/pNwqLyfl2co
 * */
func parseLine(s string) (*Command, error) {

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	m := strings.FieldsFunc(s, f)
	if l := len(m); l == 0 || l > 4 {
		return nil, errors.New("bad input")
	}

	var c Command
	c.Command = strings.ToUpper(m[0])
	c.Params = m[1:]
	return &c, nil

}

/*
 */
func showCommandHandler(session *Session, command Command) {
	/* SHOW DB
	 * SHOW DATABASES
	 * SHOW STRAT
	 * SHOW STRATEGIES
	 * SHOW INFO database
	 * SHOW SERVER */

	if len(command.Params) < 1 {
		syntaxErrorHandler(session, command)
		return
	}

	param := strings.ToUpper(command.Params[0])

	switch param {
	case "DB", "DATABASES":
		session.Connection.Writer.PrintfLine("110 database(s) present")
		session.Connection.Writer.PrintfLine(".")
		WriteCode(session, 250, "ok")
		return
	case "STRAT", "STRATEGIES":
		session.Connection.Writer.PrintfLine("111 strategie(s) present")
		session.Connection.Writer.PrintfLine(".")
		WriteCode(session, 250, "ok")
		return
	case "INFO":
		if len(command.Params) < 2 {
			syntaxErrorHandler(session, command)
			return
		}
		name := command.Params[1]
		session.Connection.Writer.PrintfLine("112 information for %s", name)
		WriteCode(session, 250, "ok")
		return
	case "SERVER":
		WriteCode(session, 114, "server information")
		WriteCode(session, 250, "ok")
		return
	}

	unknownCommandHandler(session, command)

}

func unknownCommandHandler(session *Session, command Command) {
	WriteCode(session, 500, "unknown command")
}

func handshakeHandler(session *Session) {
	session.Connection.Writer.PrintfLine("220 %s <%s> <%s>",
		"go-dictd",
		"mime",
		session.MsgId,
	)
}

/*
 */
func syntaxErrorHandler(session *Session, command Command) {
	WriteCode(session, 501, "syntax error, illegal parameters")
}

/*
 */
func quitCommandHandler(session *Session, command Command) {
	WriteCode(session, 221, "bye")
	session.Connection.Close()
}

func registerDefaultHandlers(server *Server) {
	server.RegisterHandler("SHOW", showCommandHandler)
	server.RegisterHandler("QUIT", quitCommandHandler)
}

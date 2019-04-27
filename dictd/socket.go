package dictd

import (
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

/* Session contains state data that lasts throughout the connection. */
type Session struct {
	MsgId      string
	Client     string
	Connection *textproto.Conn
	DictServer *Server
	Options    map[string]bool
}

/* Given a dict.Session and a dict.Command, route the command to the proper
 * handler, and dispatch the command. */
func handleCommand(session *Session, command *Command) {
	log.Printf("Incomming command from %s: %s", session.MsgId, command.Command)

	handler := session.DictServer.GetHandler(command)
	if handler == nil {
		unknownCommandHandler(session, *command)
	} else {
		handler(session, *command)
	}
}

/* Helper for commands to write out a text block */
func WriteTextBlock(session *Session, stream string) {
	if session.Options["MIME"] {
		session.Connection.Writer.PrintfLine(
			"Content-type: text/plain; charset=utf-8\n" +
				"Content-transfer-encoding: 8bit\n",
		)

	}

	writer := session.Connection.Writer.DotWriter()
	writer.Write([]byte(stream))
	writer.Close()
}

/* Helper for commands to write out a code line */
func WriteCode(session *Session, code int, message string) {
	session.Connection.Writer.PrintfLine("%d %s", code, message)
}

/* Helper to generate a "unique" Message ID for the client to use.
 *
 * It's not super vital, since we don't support AUTH yet. */
func generateMsgId(server *Server) string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) +
		".0@" +
		server.Name
}

/* Given a `dict.Server` and a `net.Conn`, do a bringup, and run the
 * `ReadLine` loop, dispatching commands to the correct internals. */
func Handle(server *Server, conn net.Conn) {
	proto := textproto.NewConn(conn)

	session := Session{
		MsgId:      generateMsgId(server),
		Client:     "",
		Connection: proto,
		DictServer: server,
		Options:    map[string]bool{},
	}

	session.Options["MIME"] = false /* Requiredish */

	/* Right, so we've got a connection, let's send the 220 and let the
	 * client know we're happy. */
	handshakeHandler(&session)

	for {
		line, err := proto.ReadLine()
		if err != nil {
			log.Printf("Error: %s", err)
			/* Usually an EOF */
			return
		}

		line = strings.Trim(line, " \n\r\t")
		if line == "" {
			continue
		}

		command, err := parseLine(line)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		handleCommand(&session, command)
	}
}

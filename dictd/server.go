package dictd

/* Server encapsulation.
 *
 * This contains a bundle of useful helpers, as well as a few data structures
 * to handle registered Databases and Commands. */
type Server struct {
	Name     string
	Info     string
	commands map[string]func(*Session, Command)
}

/* GetHandler returns a Command handler for the given dict.Command `command` */
func (s *Server) GetHandler(command *Command) func(*Session, Command) {
	name := command.Command

	if value, ok := s.commands[name]; ok {
		return value
	}
	return nil
}

/* Register a Command `handler` under name `name`. */
func (s *Server) RegisterHandler(
	name string,
	handler func(*Session, Command),
) {
	s.commands[name] = handler
}

/* Create a new server by name `name`. */
func NewServer(name string) Server {
	server := Server{
		Name:     name,
		Info:     "",
		commands: map[string]func(*Session, Command){},
	}

	registerDefaultHandlers(&server)
	return server
}

package main

import (
	"fmt"
	"github.com/narslan/dictd"
)

func main() {
	server := dictd.New("localhost:9999")

	server.OnNewClient(func(c *dictd.Client) {
		// new client connected
		// lets send some message
		c.Send("Hello\n")
	})
	server.OnNewMessage(func(c *dictd.Client, message string) {
		c.Send(fmt.Sprintf("Received message %s", message))
		// new message received
	})
	server.OnClientConnectionClosed(func(c *dictd.Client, err error) {
		// connection with client lost
	})

	server.Listen()
}

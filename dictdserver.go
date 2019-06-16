package main

import (
	"log"
	"net"

	"github.com/narslan/dictdserver/database"
	"pault.ag/go/dictd/dictd"
)

func main() {

	server := dictd.NewServer("dictd server")

	r := database.ReadConfig("/etc/dict/dictd.conf")
	for _, v := range r {

		server.RegisterDatabase(v, v.Name, true)

	}

	link, err := net.Listen("tcp", ":2628")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := link.Accept()
		if err != nil {
			log.Printf("Error: %s", err)
		}
		go dictd.Handle(&server, conn)
	}
}

package main

import (
	"log"
	"net"

	"github.com/narslan/dictdserver/database"
	"pault.ag/go/dictd/dictd"
)

func main() {

	server := dictd.NewServer("dictd server")

	db, err := database.NewDictdDatabase(
		"/usr/share/dictd/deu-eng.index",
		"/usr/share/dictd/deu-eng.dict.dz",
		"German English Database",
	)

	if err != nil {
		log.Fatal(err)
	}

	db2, err2 := database.NewDictdDatabase(
		"/usr/share/dictd/wn.index",
		"/usr/share/dictd/wn.dict.dz",
		"Word Net English Database",
	)

	if err2 != nil {
		log.Fatal(err2)
	}

	server.RegisterDatabase(db, "deu-eng", true)
	server.RegisterDatabase(db2, "wn", true)

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

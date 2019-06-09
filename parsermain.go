package main

import (
	"fmt"
	"github.com/narslan/dictdserver/database"
)

func main() {

	r := database.ReadConfig("/etc/dict/dictd.conf")
	for k, v := range r {
		fmt.Printf("%d %s\n", k, v)
	}
	//fmt.Printf("%T\n", r)
}

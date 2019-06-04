package main

import (
	"fmt"
	"github.com/narslan/dictdserver/database"
)

func main() {

	t1 := database.NewTree() //main tree
	t1.Root = t1.NewList(0)
	ns := t1.NewString("Newroz", 2)
	ns2 := t1.NewString("Newroz2", 5)
	ns3 := t1.NewString("Newroz3", 10)
	t1.Root.Append(ns)
	t1.Root.Append(ns2)
	t1.Root.Append(ns3)
	t2 := t1.NewList(10)

	t1.Root.Append(t2)
	ts3 := t1.NewString("Newroz4", 10)
	ts4 := t1.NewString("Newroz4", 14)
	ts5 := t1.NewMap("key", "Newroz4", 16)
	t2.Append(ts3)
	t2.Append(ts5)
	t2.Append(ts4)
	fmt.Printf("%s", t1)
}

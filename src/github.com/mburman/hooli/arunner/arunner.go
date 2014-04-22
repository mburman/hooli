package main

import (
	"flag"
	"fmt"
	"github.com/mburman/hooli/acceptor"
)

const defaultAcceptorPort = 9010

var (
	acceptorPort = flag.Int("aport", defaultAcceptorPort, "port number for accepter to listen on")
)

func printAFlags() {
	fmt.Println("Arguments...")
	fmt.Println("Acceptor Port: ", *acceptorPort)
}

func main() {
	flag.Parse()
	printAFlags()

	fmt.Println("Starting up acceptor.")
	acceptor.NewAcceptor(*acceptorPort)
	select {}
}

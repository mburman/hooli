package main

import (
	"errors"
	"flag"
	"fmt"
	//	"github.com/mburman/hooli/acceptor"
	//	"github.com/mburman/hooli/proposer"
	"strings"
)

const defaultPort = 9009

var (
	nodePort = flag.Int("port", defaultPort, "port number for accepter to listen on")
	nodeID   = flag.Uint("id", 0, "unique id to be used as lower order bits of the proposal number")
)

type ports []string

var acceptorPorts ports

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (p *ports) Set(value string) error {
	if len(*p) > 0 {
		return errors.New("interval flag already set")
	}
	for _, port := range strings.Split(value, ",") {
		*p = append(*p, port)
	}
	return nil
}

func (p *ports) String() string {
	return fmt.Sprint(*p)
}

func init() {
	flag.Var(&acceptorPorts, "ports", "comma-separated list of accepter ports")
}

func printFlags() {
	fmt.Println("Arguments...")
	fmt.Println("nodePort: ", *nodePort)
	for _, port := range acceptorPorts {
		fmt.Println("Acceptor port: ", port)
	}
	fmt.Println("nodeID: ", *nodeID)
}

func main() {
	flag.Parse()
	printFlags()

	// START UP ACCEPTOR RPC SERVER.

	// START LISTENING FOR MESSAGES FROM CLIENTS.
}

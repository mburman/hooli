package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"strings"
	"time"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type ports []string

var proposerPorts ports

var clients = make([]*rpc.Client, 0)

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
	flag.Var(&proposerPorts, "ports", "comma-separated list of proposer ports")
}

func printFlags() {
	fmt.Println("Arguments...")
	for _, port := range proposerPorts {
		fmt.Println("Proposer port: ", port)
	}
}

func main() {
	flag.Parse()

	if len(proposerPorts) < 1 {
		fmt.Println("No proposers specified")
		os.Exit(1)
	}

	// Connect to a bunch of proposers
	for _, port := range proposerPorts {
		fmt.Println("Connect to server on port: ", port)
		client, err := rpc.DialHTTP("tcp", ":"+port)
		// Keep redialing till we connect to the server.
		for err != nil {
			time.Sleep(time.Second)
			LOGE.Println("client redialing because of error: ", err)
			client, err = rpc.DialHTTP("tcp", ":"+port)
		}
		clients = append(clients, client)
	}

	// Send a bunch of random messages to random proposers
	rand.Seed(time.Now().Unix())
	done1 := make(chan int)
	done2 := make(chan int)
	done3 := make(chan int)
	go sendRandomMessagesToRandomClients(clients, 1, done1, 10)
	go sendRandomMessagesToRandomClients(clients, 0, done2, 50)
	go sendRandomMessagesToRandomClients(clients, 1, done3, 10)

	//go sendSingleMessageToAllClients(clients, 0, done1)
	<-done1
	<-done2
	<-done3
	// Wait for paxos to complete
	time.Sleep(time.Second * 5)

	// Get messages from all the servers
	messagesList := make([][]proposerrpc.Message, len(clients))
	for i := 0; i < len(clients); i++ {
		// get amd compare messages from all clients
		fmt.Println("Getting messages from client:", i)
		messagesList[i] = getMessages(clients[i])
	}

	// Verity the lengths match
	length := len(messagesList[0])
	fmt.Println("Number of messages recv by client 0:", length)
	for i := 1; i < len(messagesList); i++ {
		fmt.Printf("Number of messages recv by client %d:%d\n", i, len(messagesList[i]))
		if length != len(messagesList[i]) {
			fmt.Println("FAIL: Message lengths do not match: ", i)
			os.Exit(1)
		}
	}

	// Iterate over messages and check that consistency is maintained over
	// different acceptors
OUTER:
	for i := 0; i < length; i++ {
		message1 := &messagesList[0][i] // extract first message
		for j := 1; j < len(messagesList); j++ {
			ret := areMessagesEqual(message1, &messagesList[j][i])
			if !ret {
				fmt.Println("FAIL messages not equal")
				fmt.Println("Index: ", i)
				fmt.Println("Messages: ", message1, &messagesList[j][i])
				break OUTER
			}
		}
	}

	fmt.Println("\nPrinting Messages Recv\n")
	for i := 0; i < len(clients); i++ {
		fmt.Println("Messages recv by client: ", i)
		for j := 0; j < len(messagesList[i]); j++ {
			fmt.Println(messagesList[i][j])
		}
		fmt.Println("")
	}
}

func areMessagesEqual(message1 *proposerrpc.Message, message2 *proposerrpc.Message) bool {
	if message1.Author == message2.Author && message1.MessageText == message2.MessageText {
		return true
	}
	return false
}

func sendSingleMessageToAllClients(clients []*rpc.Client, delay int, done chan int) {
	for i := 0; i < len(clients); i++ {
		userid := rand.Int()
		fmt.Println("Sending message to client:", i)
		message := &proposerrpc.Message{
			Author:  fmt.Sprintf("%d", userid),
			MessageText: fmt.Sprintf("%d", i),
		}
		sendMessage(i, message)
		time.Sleep(time.Second * time.Duration(delay)) // wait for paxos to complete
	}
	close(done)
}

func sendRandomMessagesToRandomClients(clients []*rpc.Client, delay int, done chan int, numMessages int) {
	for i := 0; i < numMessages; i++ {
		proposerToContact := rand.Intn(len(clients)) // pick a random client
		userid := rand.Int()
		fmt.Println("Sending message to client:", proposerToContact)
		message := &proposerrpc.Message{
			Author:  fmt.Sprintf("%d", userid),
			MessageText: fmt.Sprintf("%d", i),
		}
		sendMessage(proposerToContact, message)
		time.Sleep(time.Second * time.Duration(delay)) // wait for paxos to complete
	}
	close(done)
}

func sendMessage(clientIndex int, message *proposerrpc.Message) {
	request := &proposerrpc.PostMessageArgs{
		Message: *message,
	}
	client := clients[clientIndex]
	var reply proposerrpc.PostMessageReply
	err := client.Call("ProposerObj.PostMessageRPC", request, &reply)
	if err != nil {
		LOGE.Println("rpc error: ", err, clientIndex, len(proposerPorts), len(clients))

		// Redial if disconnected.
		// TODO: stop redialing after a bit...
		client, err := rpc.DialHTTP("tcp", ":"+proposerPorts[clientIndex])
		if err != nil {
			LOGE.Println("error in redial:", err)
			time.Sleep(time.Second * 1)
		} else {
			clients[clientIndex] = client
			sendMessage(clientIndex, message)
		}
	}
}

func getMessages(client *rpc.Client) []proposerrpc.Message {
	request := &proposerrpc.GetMessagesArgs{}

	var reply proposerrpc.GetMessagesReply
	err := client.Call("ProposerObj.GetMessages", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}

	return reply.Messages
}

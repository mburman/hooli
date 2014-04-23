package main

import (
	"flag"
	"fmt"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"time"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

var (
	masterServerHostPort = flag.String("host:port", "", "host:port to contact the server on.")
)

func main() {
	flag.Parse()
	fmt.Println("Contacting server on port: ", *masterServerHostPort)

	client, err := rpc.DialHTTP("tcp", *masterServerHostPort)
	// Keep redialing till we connect to the server.
	for err != nil {
		time.Sleep(time.Second)
		LOGE.Println("client redialing because of error: ", err)
		client, err = rpc.DialHTTP("tcp", *masterServerHostPort)
	}
	LOGE.Println("client dialed successfully")

	rand.Seed(time.Now().Unix())
	userid := rand.Int()
	for i := 0; i <= 10; i++ {
		fmt.Println("ITERATION: ", i)
		message := &proposerrpc.Message{
			Userid:  fmt.Sprintf("%d", userid),
			Message: fmt.Sprintf("%d", i),
		}
		sendMessage(client, message)
		time.Sleep(time.Second * 2) // wait for paxos to complete
		getMessages(client)
	}
}

func sendMessage(client *rpc.Client, message *proposerrpc.Message) {
	request := &proposerrpc.PostMessageArgs{
		Message: *message,
	}

	var reply proposerrpc.PostMessageReply
	err := client.Call("ProposerObj.PostMessage", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
}

func getMessages(client *rpc.Client) {
	request := &proposerrpc.GetMessagesArgs{}

	var reply proposerrpc.GetMessagesReply
	err := client.Call("ProposerObj.GetMessages", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}

	fmt.Println("List of messages:")
	for _, message := range reply.Messages {
		fmt.Printf("%+v\n", message)
	}
}

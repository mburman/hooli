package main

import (
	"flag"
	"fmt"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
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

	request := &proposerrpc.PostMessageArgs{
		Message: proposerrpc.Message{
			Userid:  "yoloid",
			Message: "yolo swag brah",
		},
	}

	var reply proposerrpc.PostMessageReply
	err = client.Call("ProposerObj.PostMessage", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}

	fmt.Println("Done posting.")
}

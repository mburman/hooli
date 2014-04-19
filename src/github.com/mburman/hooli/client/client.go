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

	// Keep redialing till we connect to the server.
	client, err := rpc.DialHTTP("tcp", *masterServerHostPort)
	if err != nil {
		LOGE.Println("redialing:", err)
		dialTicker := time.NewTicker(time.Second)
	DialLoop:
		for {
			select {
			case <-dialTicker.C:
				client, err = rpc.DialHTTP("tcp", *masterServerHostPort)
				if err == nil {
					break DialLoop
				}
			}
		}
	}

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

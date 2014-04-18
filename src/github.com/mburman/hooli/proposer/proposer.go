package proposer

import (
	"fmt"
	//"github.com/mburman/hooli/acceptor"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type proposerObj struct {
}

// port: port to start the proposerObj on.
func NewProposer(port int) *proposerObj {
	var a proposerObj
	setupRPC(&a, port)
	return &a
}

// Client calls this to post a message.
func (a *proposerObj) PostMessage(args *proposerrpc.PostMessageArgs, reply *proposerrpc.PostMessageReply) error {
	return nil
}

func (a *proposerObj) GetMessages(args *proposerrpc.GetMessagesArgs, reply *proposerrpc.GetMessagesReply) error {
	return nil
}

func setupRPC(a *proposerObj, port int) {
	rpc.RegisterName("ProposerObj", a)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		LOGE.Println("listen error:", e)
	}
	go http.Serve(l, nil)
}

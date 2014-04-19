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
	port          int
	acceptorPorts []int
	messageQueue  chan *proposerrpc.Message // messages to be handled
	messages      []proposerrpc.Message     // list of messages
}

// port: port for proposer to listen to client requests on.
// acceptorPorts: ports contact acceptors on.
func NewProposer(port int, acceptorPorts []int) *proposerObj {
	var p proposerObj
	p.port = port
	p.acceptorPorts = acceptorPorts
	p.messageQueue = make(chan *proposerrpc.Message, 100)

	setupRPC(&p, port)
	go processMessages(&p) // start processing incoming messa
	return &p
}

// Client calls this to post a message.
func (p *proposerObj) PostMessage(args *proposerrpc.PostMessageArgs, reply *proposerrpc.PostMessageReply) error {
	// Promise to handle. Don't block.
	go handleMessage(p, &args.Message)
	return nil
}

func (p *proposerObj) GetMessages(args *proposerrpc.GetMessagesArgs, reply *proposerrpc.GetMessagesReply) error {
	// TODO: algorithm for figuring out nearest.
	reply.Messages = p.messages
	return nil
}

func handleMessage(p *proposerObj, message *proposerrpc.Message) {
	p.messageQueue <- message
}

// Continuously reads messages from queue and Paxos' them
func processMessages(p *proposerObj) {
	for {
		message := <-p.messageQueue
		message = message
		// TODO: paxos
	}
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

package proposer

import (
	"fmt"
	//"github.com/mburman/hooli/acceptor"
	"github.com/mburman/hooli/rpc/acceptorrpc"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type proposerObj struct {
	port          int
	acceptorPorts []int
	messageQueue  chan *proposerrpc.Message // messages to be handled
	messages      []proposerrpc.Message     // list of messages
	acceptorList  []*rpc.Client
}

// port: port for proposer to listen to client requests on.
// acceptorPorts: ports contact acceptors on.
func NewProposer(port int, acceptorPorts []int) *proposerObj {
	var p proposerObj
	p.port = port
	p.acceptorPorts = acceptorPorts
	p.messageQueue = make(chan *proposerrpc.Message, 100)
	p.acceptorList = make([]*rpc.Client, 0)

	setupRPC(&p, port)
	connectToAcceptors(&p)
	go processMessages(&p) // start processing incoming messa
	return &p
}

// Client calls this to post a message.
func (p *proposerObj) PostMessage(args *proposerrpc.PostMessageArgs, reply *proposerrpc.PostMessageReply) error {
	// Promise to handle. Don't block.
	fmt.Printf("Received Message:  %+v\n", args.Message)
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

func connectToAcceptors(p *proposerObj) {
	// Connect to acceptors.
	for acceptorPort := range p.acceptorPorts {
		// Keep redialing till we connect to the server.
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", acceptorPort))
		if err != nil {
			LOGE.Println("redialing:", err)
			dialTicker := time.NewTicker(time.Second)
		DialLoop:
			for {
				select {
				case <-dialTicker.C:
					client, err = rpc.DialHTTP("tcp", fmt.Sprintf(":%d", acceptorPort))
					if err == nil {
						break DialLoop
					}
				}
			}
		}
		p.acceptorList = append(p.acceptorList, client)
	}
}

func chooseProposalNumber(p *proposerObj) int {
	// TODO:
	return 0
}

func sendPrepare(p *proposerObj, client *rpc.Client, proposalNumber int) acceptorrpc.PrepareReply {
	// TODO:
	return acceptorrpc.PrepareReply{}
}

func sendAccept(p *proposerObj, client *rpc.Client, proposalNumber int,
	acceptedMessage *proposerrpc.Message) acceptorrpc.AcceptReply {
	// TODO:
	return acceptorrpc.AcceptReply{}
}

// Continuously reads messages from queue and Paxos' them
func processMessages(p *proposerObj) {
	for {
		message := <-p.messageQueue
		for {
			proposalNumber := chooseProposalNumber(p)

			// Send prepares.
			acceptedMessage := message
			for _, a := range p.acceptorList {
				// Send prepare message to all the acceptors
				prepareReply := sendPrepare(p, a, proposalNumber)
				if prepareReply.AcceptedMessage.Message != "" {
					acceptedMessage = &prepareReply.AcceptedMessage
					break
				}
			}

			// Send accepts.
			for _, a := range p.acceptorList {
				acceptReply := sendAccept(p, a, proposalNumber, acceptedMessage)
				if acceptReply.MinProposal > proposalNumber {
					break
				}
			}
		}
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

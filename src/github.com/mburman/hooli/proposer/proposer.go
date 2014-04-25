package proposer

import (
	"fmt"
	//"github.com/mburman/hooli/acceptor"
	"github.com/mburman/hooli/rpc/acceptorrpc"
	"github.com/mburman/hooli/rpc/proposerrpc"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

var LOGE = log.New(os.Stderr, "ERROR ", log.Lmicroseconds|log.Lshortfile)
var LOGV = log.New(ioutil.Discard, "VERBOSE ", log.Lmicroseconds|log.Lshortfile)

type proposerObj struct {
	port              int
	acceptorPorts     []string
	messageQueue      chan *proposerrpc.Message // messages to be handled
	messages          []proposerrpc.Message     // list of messages
	acceptorList      []*rpc.Client
	maxProposalNumber int // max proposal number server has seen
	id                int // proposer id
}

// port: port for proposer to listen to client requests on.
// acceptorPorts: ports contact acceptors on.
func NewProposer(port int, acceptorPorts []string) *proposerObj {
	var p proposerObj
	p.port = port
	p.acceptorPorts = acceptorPorts
	p.messageQueue = make(chan *proposerrpc.Message, 100)
	p.acceptorList = make([]*rpc.Client, 0)

	rand.Seed(time.Now().Unix())
	p.id = rand.Intn(100)               // Random server id.
	p.maxProposalNumber = rand.Intn(10) // Randomize initial round number.

	setupRPC(&p, port)
	connectToAcceptors(&p)
	go processMessages(&p) // start processing incoming messa
	return &p
}

// Client calls this to post a message.
func (p *proposerObj) PostMessage(args *proposerrpc.PostMessageArgs, reply *proposerrpc.PostMessageReply) error {
	// Promise to handle. Don't block.
	fmt.Printf("Received Message to post:  %+v\n", args.Message)
	go handleMessage(p, &args.Message)
	return nil
}

func (p *proposerObj) GetMessages(args *proposerrpc.GetMessagesArgs, reply *proposerrpc.GetMessagesReply) error {
	LOGV.Println("Getting messages")
	// Get messages from a storage server (acceptor)
	if len(p.acceptorList) < 1 {
		// panic return error
	}

	client := p.acceptorList[0]
	request := &acceptorrpc.GetMessagesArgs{
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Radius:    args.Radius,
	}
	var messagesReply acceptorrpc.GetMessagesReply
	err := client.Call("AcceptorObj.GetMessages", request, &messagesReply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	reply.Messages = messagesReply.Messages
	return nil
}

func handleMessage(p *proposerObj, message *proposerrpc.Message) {
	p.messageQueue <- message
}

func connectToAcceptors(p *proposerObj) {
	// Connect to acceptors.
	for _, acceptorPort := range p.acceptorPorts {
		// Keep redialing till we connect to the server.
		client, err := rpc.DialHTTP("tcp", ":"+acceptorPort)
		if err != nil {
			LOGE.Println("redialing:", err)
			dialTicker := time.NewTicker(time.Second)
		DialLoop:
			for {
				select {
				case <-dialTicker.C:
					client, err = rpc.DialHTTP("tcp", ":"+acceptorPort)
					if err == nil {
						break DialLoop
					}
				}
			}
		}
		p.acceptorList = append(p.acceptorList, client)
	}
}

func generateProposal(p *proposerObj) *acceptorrpc.Proposal {
	p.maxProposalNumber++
	return &acceptorrpc.Proposal{
		Number: p.maxProposalNumber,
		ID:     p.id,
	}
}

func sendPrepare(p *proposerObj, client *rpc.Client, proposal *acceptorrpc.Proposal) acceptorrpc.PrepareReply {
	request := &acceptorrpc.PrepareArgs{
		Proposal: *proposal,
	}
	var reply acceptorrpc.PrepareReply
	err := client.Call("AcceptorObj.Prepare", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

func sendAccept(p *proposerObj, client *rpc.Client, proposal *acceptorrpc.Proposal,
	acceptedMessage *proposerrpc.Message) acceptorrpc.AcceptReply {
	request := &acceptorrpc.AcceptArgs{
		Proposal:        *proposal,
		ProposalMessage: *acceptedMessage,
	}
	var reply acceptorrpc.AcceptReply
	err := client.Call("AcceptorObj.Accept", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

func sendCommit(p *proposerObj, client *rpc.Client, message *proposerrpc.Message, index int) acceptorrpc.CommitReply {
	request := &acceptorrpc.CommitArgs{
		Message: *message,
		Index:   index,
	}
	var reply acceptorrpc.CommitReply
	err := client.Call("AcceptorObj.Commit", request, &reply)
	if err != nil {
		LOGE.Println("rpc error:", err)
	}
	return reply
}

// Continuously reads messages from queue and Paxos' them
func processMessages(p *proposerObj) {
	for {
		message := <-p.messageQueue
		for {
			delayMills := rand.Intn(75) // Random delay to avoid conflicts.
			time.Sleep(time.Millisecond * time.Duration(delayMills))
			proposal := generateProposal(p)

			// Send prepares. TODO: needs to be done async since nodes might go down.
			acceptCount := 0
			highestCancelProposalNumber := -1
			highestIndex := 0
			for _, a := range p.acceptorList {
				// Send prepare message to all the acceptors
				prepareReply := sendPrepare(p, a, proposal)
				if prepareReply.Status == acceptorrpc.OK {
					if highestIndex < prepareReply.Index {
						highestIndex = prepareReply.Index
					}
					acceptCount++
				} else if prepareReply.Status == acceptorrpc.PREV_ACCEPTED {
					// IGNORE...
				} else if prepareReply.Status == acceptorrpc.CANCEL {
					if prepareReply.AcceptedProposalNumber >= highestCancelProposalNumber {
						highestCancelProposalNumber = prepareReply.AcceptedProposalNumber
					}
				}
			}

			// If a majority have not accepted - this is not the leader
			// Try again.
			if acceptCount <= len(p.acceptorList)/2 {
				p.maxProposalNumber = highestCancelProposalNumber
				continue
			}

			fmt.Println("Sending accepts")
			// LEADER. Send accepts.
			acceptCount = 0
			highestCancelProposalNumber = -1
			for _, a := range p.acceptorList {
				acceptReply := sendAccept(p, a, proposal, message)
				if acceptReply.Status == acceptorrpc.OK {
					acceptCount++
				} else if acceptReply.Status == acceptorrpc.CANCEL {
					if acceptReply.MinProposalNumber >= highestCancelProposalNumber {
						highestCancelProposalNumber = acceptReply.MinProposalNumber
					}
				} else {
					// ASSERT FALSE should never happen
				}
			}

			if acceptCount <= len(p.acceptorList)/2 {
				p.maxProposalNumber = highestCancelProposalNumber
				continue
			}

			fmt.Println("COMMITTING")
			// Value has been chosen.
			for _, a := range p.acceptorList {
				sendCommit(p, a, message, highestIndex)
			}

			break // get a new message.
		}
	}
}

func setupRPC(a *proposerObj, port int) {
	fmt.Println("Proposer rpc:", port)
	rpc.RegisterName("ProposerObj", a)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		LOGE.Println("listen error:", e)
	}
	go http.Serve(l, nil)
}

package acceptor

import (
	"fmt"
	"github.com/mburman/hooli/rpc/acceptorrpc"
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

type acceptorObj struct {
	minProposal     *acceptorrpc.Proposal
	acceptedMessage *proposerrpc.Message
	messages        []proposerrpc.Message
}

// port: port to start the acceptorObj on.
func NewAcceptor(port int) *acceptorObj {
	var a acceptorObj
	a.minProposal = &acceptorrpc.Proposal{Number: -1, ID: -1}
	a.acceptedMessage = nil
	a.messages = make([]proposerrpc.Message, 0)

	setupRPC(&a, port)
	return &a
}

func (a *acceptorObj) Prepare(args *acceptorrpc.PrepareArgs, reply *acceptorrpc.PrepareReply) error {
	fmt.Println("Received Prepare")
	if args.Proposal.Number < a.minProposal.Number {
		fmt.Println("CANCEL")
		reply.Status = acceptorrpc.CANCEL
	} else if args.Proposal.Number == a.minProposal.Number && args.Proposal.ID < a.minProposal.ID {
		fmt.Println("CANCEL EQUAL")
		reply.Status = acceptorrpc.CANCEL
	} else {
		if a.minProposal.Number != -1 {
			// Something was previously accepted.
			fmt.Println("PREV ACCEPTED")
			reply.Status = acceptorrpc.PREV_ACCEPTED
		} else {
			fmt.Println("OK")
			reply.Status = acceptorrpc.OK
		}
	}

	reply.AcceptedProposalNumber = a.minProposal.Number
	if a.acceptedMessage != nil {
		reply.AcceptedMessage = *a.acceptedMessage
	}
	return nil
}

func (a *acceptorObj) Accept(args *acceptorrpc.AcceptArgs, reply *acceptorrpc.AcceptReply) error {
	fmt.Println("Received Accept.")
	if args.Proposal.Number < a.minProposal.Number {
		reply.Status = acceptorrpc.CANCEL
	} else if args.Proposal.Number == a.minProposal.Number && args.Proposal.ID < a.minProposal.ID {
		reply.Status = acceptorrpc.CANCEL
	} else {
		reply.Status = acceptorrpc.OK
		a.acceptedMessage = &args.ProposalMessage
		a.minProposal = &args.Proposal
	}

	reply.MinProposalNumber = a.minProposal.Number
	return nil
}

func (a *acceptorObj) Commit(args *acceptorrpc.CommitArgs, reply *acceptorrpc.CommitReply) error {
	fmt.Println("Committed message:", args.Message)
	a.messages = append(a.messages, args.Message)
	return nil
}

func (a *acceptorObj) GetMessages(args *acceptorrpc.GetMessagesArgs, reply *acceptorrpc.GetMessagesReply) error {
	fmt.Println("Getting ", len(a.messages), " messages")
	reply.Messages = a.messages
	return nil
}

func setupRPC(a *acceptorObj, port int) {
	fmt.Println("Acceptor rpc:", port)
	rpc.RegisterName("AcceptorObj", a)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		LOGE.Println("listen error:", e)
	}
	go http.Serve(l, nil)
}

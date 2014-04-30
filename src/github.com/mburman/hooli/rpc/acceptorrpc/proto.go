// Contains accept RPCs

package acceptorrpc

import (
	"github.com/mburman/hooli/rpc/proposerrpc"
)

type Status int

const (
	OK            Status = iota + 1 // The RPC was a success.
	PREV_ACCEPTED                   // Something else was previously accepted.
	CANCEL
	ALREADY_FILLED
)

type Proposal struct {
	Number int
	ID     int // unique server id to break clashing numbers
}

type PrepareArgs struct {
	Proposal Proposal
}

type PrepareReply struct {
	AcceptedProposalNumber int
	AcceptedMessage        proposerrpc.Message
	Index                  int
	Status                 Status
}

type AcceptArgs struct {
	Proposal        Proposal
	ProposalMessage proposerrpc.Message
	Index           int // send index where we want to commit
}

type AcceptReply struct {
	MinProposalNumber int
	Status            Status
	Message           proposerrpc.Message
}

type CommitArgs struct {
	Message proposerrpc.Message
	Index   int
}

type CommitReply struct {
}

type GetMessagesArgs struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type GetMessagesReply struct {
	Messages []proposerrpc.Message
}
